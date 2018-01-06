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

	//bfCode = optimizeLoops(bfCode)
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

//func optimizeLoops(ops []Op) []Op {
//	for i := 0; i < len(ops); i++ {
//		op := ops[i]
//		if op == StartLoop {
//			end := beunSearch(ops, i)
//			simpleLoop, err := simplifyLoop(ops[i:end])
//			if err != nil {
//				continue
//			}
//			if len(simpleLoop) > end-i {
//				logE.Printf("What good does")
//			}
//			//logE.Printf("Current operation is GOOD: %v\n", ops[end] == EndLoop)
//			//logE.Printf("Correct one is -1: %v, +1: %v", ops[end-1] == EndLoop, ops[end+1] == EndLoop)
//			//panic("aaaaa")
//		}
//	}
//	return ops
//}
//func simplifyLoop(ops []Op) ([]Op, error) {
//	index := 0
//	changeToBaseNum := 0
//
//	change := make([]uint8, BUFSIZE)
//	for i := 0; i < len(ops); i++ {
//		op := ops[i]
//		switch op {
//		case Print, Input:
//			return ops, errors.New("loop contains IO, cannot optimize")
//		case StartLoop, EndLoop:
//			return ops, errors.New("Cannot optimize loops with sub-loops yet")
//		case IndexDec:
//			index--
//		case IndexDecArg:
//			index -= int(ops[i+1])
//			i++
//		case IndexInc:
//			index++
//		case IndexIncArg:
//			index += int(ops[i+1])
//			i++
//		case DataDec:
//			if index == 0 {
//				changeToBaseNum--
//			}
//		case DataDecArg:
//			if index == 0 {
//				changeToBaseNum -= int(ops[i+1])
//				i++
//			}
//		case DataInc:
//			if index == 0 {
//				changeToBaseNum++
//			}
//		case DataIncArg:
//			if index == 0 {
//				changeToBaseNum += int(ops[i+1])
//				i++
//			}
//		}
//	}
//	if index == 0 && changeToBaseNum == -1 {
//		// This loop can be optimized! We can just apply it `n` times!
//
//	}
//}
//
////
//func beunSearch(s []Op, i int) int {
//	count := 1
//	for i += 1; count > 0; i++ {
//		if s[i] == StartLoop {
//			count++
//		} else if s[i] == EndLoop {
//			count--
//		}
//	}
//	return i - 1
//}

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

func optimize(b *[]byte) ([]Op, map[uint]uint) {
	s := *b
	const valid = "><+-[].,"
	var optimized []Op

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

func toOp(c uint8) Op {
	switch c {
	case '<':
		return IndexDec
	case '>':
		return IndexInc
	case '+':
		return DataInc
	case '-':
		return DataDec
	case '.':
		return Print
	case ',':
		return Input
	default:
		logE.Printf("This is not a valid Op: %d!", c)
		return IndexDec
	}
}

func toOpWithArg(c uint8) Op {
	switch c {
	case '<':
		return IndexDecArg
	case '>':
		return IndexIncArg
	case '+':
		return DataIncArg
	case '-':
		return DataDecArg
	default:
		logE.Printf("This is not a valid Op: %d!", c)
		return IndexDec
	}
}

func isOpWithArg(op Op) bool {
	return op >= DataDecArg && op <= IndexIncArg
}

func bfExecute(ptr *[]Op, state *State) {
	bf := *ptr
	switch op := bf[state.instr]; op {
	case IndexDec:
		state.IndexDec(1)
	case IndexDecArg:
		n := uint(bf[state.instr+1])
		state.IndexDec(n)
		state.instr++
	case IndexInc:
		state.IndexInc(1)
	case IndexIncArg:
		n := uint(bf[state.instr+1])
		state.IndexInc(n)
		state.instr++
	case DataDec:
		state.DataDec(1)
	case DataDecArg:
		n := uint(bf[state.instr+1])
		state.DataDec(n)
		state.instr++
	case DataInc:
		state.DataInc(1)
	case DataIncArg:
		n := uint(bf[state.instr+1])
		state.DataInc(n)
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
	case Plus:
		state.Plus()
	case Minus:
		state.Minus()
	case Mult:
		state.Mult()

	default:
		logE.Printf("This cannot happen: op = %d", op)
	}
	//fmt.Printf("%op ", bf[state.instr])
	state.instr++
}
