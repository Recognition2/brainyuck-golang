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

type Op uint8

// Define all operations
const (
	// Traditional ops
	DataDec Op = iota
	DataInc
	IndexDec
	IndexInc
	Print
	Input
	StartLoop // Special
	EndLoop

	// Operations that take an argument
	DataDecArg
	DataIncArg
	IndexDecArg
	IndexIncArg

	// Advanced ops
	Zero
	Move
	Copy
	Plus
	Minus
	Mult
	Exp
	Divide
)

const BUFSIZE = 30000

var (
	logE = log.New(os.Stderr, "[ERRO] ", log.Ldate+log.Ltime+log.Ltime+log.Lshortfile)
	logW = log.New(os.Stdout, "[WARN] ", log.Ldate+log.Ltime)
	logI = log.New(os.Stdout, "[INFO] ", log.Ldate+log.Ltime)
)

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
	bfCode, jumpfwd := optimize(rawBF)
	state.jumpFwd = jumpfwd

	initTime := time.Since(startTime)
	fmt.Println("Done optimizing, running...")
	// Run all instructions
	for state.instr < uint(len(bfCode)) {
		bfExecute(bfCode, &state)
	}
	fmt.Printf("\n%s\n", state.output)
	runTime := time.Since(startTime)
	fmt.Printf("Optimizing took %s.\nTotal took %s\n", initTime, runTime)
	state.printStats()
	return
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

func optimize(s []byte) ([]Op, map[uint]uint) {
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

func bfExecute(bf []Op, state *State) {
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
		state.StartLoop(bf)
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
