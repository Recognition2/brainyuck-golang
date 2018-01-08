package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const BUFSIZE = 30000

var (
	logE = log.New(os.Stderr, "[ERRO] ", log.Ldate+log.Ltime+log.Ltime+log.Lshortfile)
	//logW = log.New(os.Stdout, "[WARN] ", log.Ldate+log.Ltime)
	//logI = log.New(os.Stdout, "[INFO] ", log.Ldate+log.Ltime)
)

const statistics = false

func main() {
	fmt.Println("Running")
	startTime := time.Now()

	var filename = flag.String("filename", "test.bf", "Path to file containing BrainYuck code")
	flag.Parse()

	rawBF, err := ioutil.ReadFile(*filename)
	if err != nil {
		logE.Println(err)
	}

	var state = GenState()

	// Build map with forward jump addresses
	bfCode, jumpfwd := optimize(&rawBF)
	state.jumpFwd = jumpfwd

	bfCode = optimizeLoops(bfCode, jumpfwd)
	fmt.Println("Done optimizing, running...")
	initTime := time.Since(startTime)

	// Run all instructions
	for state.instr < uint(len(bfCode)) {
		bfExecute(&bfCode, &state)
	}
	// Print output
	fmt.Printf("\n%s\n", state.output)

	// Timing stuffs
	runTime := time.Since(startTime)
	fmt.Printf("Optimizing took %s.\nTotal took %s\n", initTime, runTime)

	// Function call statistics
	//if statistics {
	//	state.printStats()
	//}
	return
}

func optimizeLoops(ops []Op, skipLoopMap map[uint]uint) []Op {
	var i uint
	for i = 0; i < uint(len(ops)); i++ {
		op := ops[i]
		if op == StartLoop {
			end := skipLoopMap[i]
			content := ops[i:end]
			whatLoop, whatHappensMap := analyseLoop(content)
			var newOps = make([]Op, 0)
			switch whatLoop {
			case NoLoop:
				continue
			case ZeroIndexLoop:
				newOps = optimizeZeroIndex(whatHappensMap)
			default:
			}
			//logE.Printf("Current operation is GOOD: %v\n", ops[end] == EndLoop)
			//logE.Printf("Correct one is -1: %v, +1: %v", ops[end-1] == EndLoop, ops[end+1] == EndLoop)
			//panic("aaaaa")

			// If the loop was optimized, replace the old code with the new code.
			ops = append(append(ops[:i], newOps...), ops[end+1:]...)
		}
	}
	return ops
}

func analyseLoop(ops []Op) (loop, map[int]int) {
	index := 0
	change := make(map[int]int, BUFSIZE)

	for i := 0; i < len(ops); i++ {
		op := ops[i]
		switch op {
		case Print, Input:
			return NoLoop, nil
		case StartLoop, EndLoop:
			return NoLoop, nil
		case IndexDec:
			index--
		case IndexDecArg:
			index -= int(ops[i+1])
			i++
		case IndexInc:
			index++
		case IndexIncArg:
			index += int(ops[i+1])
			i++
		case DataDec:
			change[index]--
		case DataDecArg:
			change[index] -= int(ops[i+1])
		case DataInc:
			change[index]++
		case DataIncArg:
			change[index] += int(ops[i+1])

		}
	}
	if index == 0 && change[0] == -1 {
		// This loop can be optimized! We can just apply it `n` times!
		// The indices will not change during the loop
		return ZeroIndexLoop, change
	}
	return NoLoop, nil
}

func beunSearch(s []Op, i int) int {
	count := 1
	for i += 1; count > 0; i++ {
		if s[i] == StartLoop {
			count++
		} else if s[i] == EndLoop {
			count--
		}
	}
	return i - 1
}

func startsWith(s []byte, i int, substr string) bool {
	subByte := []byte(substr)
	// Iterate over the whole substring
	for l := 0; l < len(subByte); l++ {
		// If out of bounds of beginning of string, or they are not equal
		if i+l > len(s) || s[i+l] != subByte[l] {
			return false
		}
	}
	return true
}

func optimize(b *[]byte) ([]Routine, map[uint]uint) {
	s := *b
	const valid = "><+-[].,"
	var optimized []Routine

	skipLoopMap := make(map[uint]uint, STACKSIZE) // Build map for skipping loops
	stack := make(Stack, 0)

	for i := 0; i < len(s); i++ {
		c := s[i]

		if !strings.Contains(valid, string(c)) {
			continue // Illegal characters should be ignored. This is mandatory to allow our intermediate
			// representation
		}

		switch {
		// Try to optimize `[-]` and `[+]` away
		case startsWith(s, i, "[-]"), startsWith(s, i, "[+]"):
			optimized = append(optimized, Zero)
			i += 2

			// Try to optimize `[->+<]` away, which evaluates [a b] to [0 a+b]
		case startsWith(s, i, "[->+<]"):
			optimized = append(optimized, Plus)
			i += 5

			// Subtraction, evaluates [x, y] to [x-y, 0]
		case startsWith(s, i, ">[-<->]<"):
			optimized = append(optimized, Minus)
			i += 7

			// Copy a value
		case startsWith(s, i, "[->+>+<<]>>[-<<+>>]<<"):
			optimized = append(optimized, Copy)
			i += 20

			// Exponentiation
		case startsWith(s, i, ">>+<[->[-<<[->>>+>+<<<<]>>>>[-<<<<+>>>>]<<]>[-<+>]<<]<"):
			optimized = append(optimized, Exp)
			i += 53

		case startsWith(s, i, "[>[->+>+<<]>[-<<-[>]>>>[<[-<->]<[>]>>[[-]>>+<]>-<]<<]>>>+<<[-<<+>>]<<<]>>>>>[-<<<<<+>>>>>]<<<<<"):
			optimized = append(optimized, Divide)
			i += 94

		case c == '[':
			// Try to optimize `[->+<]` away, which evaluates [a b] to [0 a+b]

			// Traditional `[` operator: start a loop.
			stack.Push(uint(len(optimized))) // Store "Begin" address
			optimized = append(optimized, StartLoop)

		case c == ']':
			skipLoopMap[stack.Pop()] = uint(len(optimized)) // map begin address to end address
			optimized = append(optimized, EndLoop)

		case c == '>', c == '<', c == '-', c == '+':
			// Try to find matching ones afterwards and collapse that
			var count uint8 = 1
			for s[i+int(count)] == c && count < 255 {
				count++
			}
			// Count now contains the amount of chars equal to current character
			if count == 1 {
				optimized = append(optimized, toOp(c))
			} else {
				optimized = append(optimized, toOpWithArg(c))
				optimized = append(optimized, Op(count)) // This is not an actual Op.
				i += int(count) - 1                      // Skip over all iterations
				// -1 because the loop adds one
			}

		default:
			optimized = append(optimized, toOp(c))
		}
	}
	return optimized, skipLoopMap
}

func bfExecute(ptr *[]Op, state *State) {

	bf := *ptr
	switch op := bf[state.instr]; op {
	case IndexDec:
		state.IndexDec(1)
	case IndexInc:
		state.IndexInc(1)
	case DataDec:
		state.DataDec(1, 0)
	case DataInc:
		state.DataInc(1, 0)

	case IndexDecArg:
		state.IndexDec(uint(bf[state.instr+1]))
		state.instr++
	case IndexIncArg:
		state.IndexInc(uint(bf[state.instr+1]))
		state.instr++
	case DataDecArg:
		state.DataDec(uint(bf[state.instr+1]), 0)
		state.instr++
	case DataIncArg:
		state.DataInc(uint(bf[state.instr+1]), 0)
		state.instr++

	case DataDecOffset:
		state.DataDec(1, state.offset)
	case DataIncOffset:
		state.DataInc(1, state.offset)
	case DataDecArgOffset:
		state.DataDec(uint(bf[state.instr+1]), state.offset)
		state.instr++
	case DataIncArgOffset:
		state.DataInc(uint(bf[state.instr+1]), state.offset)
		state.instr++

	case Print:
		state.Print()
	case Input:
		state.Input()
	case StartLoop:
		state.StartLoop()
	case EndLoop:
		state.EndLoop()

		// Special operations
	case Zero:
		state.Zero()
	case Offset:
		state.Offset(int8(bf[state.instr+1]))
		state.instr++
	case Plus:
		state.Plus()
	case Minus:
		state.Minus()
	case Mult:
		state.Mult()
	case Copy:
		state.Copy()
	case Exp:
		state.Exp()
	case Divide:
		state.Divide()

	default:
		logE.Printf("This cannot happen: op = %d", op)
	}
	//fmt.Printf("%op ", bf[state.instr])
	state.instr++
}
