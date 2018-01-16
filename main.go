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
	startTime := time.Now()
	fmt.Println("Running")

	var filename = flag.String("filename", "test.bf", "Path to file containing BrainYuck code")
	flag.Parse()

	rawBF, err := ioutil.ReadFile(*filename)
	if err != nil {
		logE.Println(err)
	}

	var state = GenState()

	// Build map with forward jump addresses
	ops := translate(rawBF)

	//ops = optimizeLoops(bfCode, jumpfwd)
	fmt.Println("Done optimizing, running...")
	initTime := time.Since(startTime)

	// Run all instructions
	for _, op := range ops {
		op.execute(&state)
	}

	// Print output
	//fmt.Printf("\n%s\n", state.output)

	// Timing stuffs
	runTime := time.Since(startTime)
	fmt.Printf("Optimizing took %s.\nTotal took %s\n", initTime, runTime)

	// Function call statistics
	if statistics {
		state.printStats()
	}
	return
}

//func optimizeLoops(ops []Op, skipLoopMap map[uint]uint) []Op {
//	var i uint
//	for i = 0; i < uint(len(ops)); i++ {
//		op := ops[i]
//		if op == StartLoop {
//			end := skipLoopMap[i]
//			content := ops[i:end]
//			whatLoop, whatHappensMap := analyseLoop(content)
//			var newOps = make([]Op, 0)
//			switch whatLoop {
//			case NoLoop:
//				continue
//			case ZeroIndexLoop:
//				newOps = optimizeZeroIndex(whatHappensMap)
//			default:
//			}
//			//logE.Printf("Current operation is GOOD: %v\n", ops[end] == EndLoop)
//			//logE.Printf("Correct one is -1: %v, +1: %v", ops[end-1] == EndLoop, ops[end+1] == EndLoop)
//			//panic("aaaaa")
//
//			// If the loop was optimized, replace the old code with the new code.
//			ops = append(append(ops[:i], newOps...), ops[end+1:]...)
//		}
//	}
//	return ops
//}

//func analyseLoop(ops []Op) (Loop, map[int]int) {
//	index := 0
//	change := make(map[int]int, BUFSIZE)
//
//	for i := 0; i < len(ops); i++ {
//		op := ops[i]
//		switch op {
//		case Print, Input:
//			return NoLoop, nil
//		case StartLoop, EndLoop:
//			return NoLoop, nil
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
//			change[index]--
//		case DataDecArg:
//			change[index] -= int(ops[i+1])
//		case DataInc:
//			change[index]++
//		case DataIncArg:
//			change[index] += int(ops[i+1])
//
//		}
//	}
//	if index == 0 && change[0] == -1 {
//		// This loop can be optimized! We can just apply it `n` times!
//		// The indices will not change during the loop
//		return ZeroIndexLoop, change
//	}
//	return NoLoop, nil
//}

func beunSearch(s []byte, i int) int {
	count := 1
	for i += 1; count > 0; i++ {
		if i >= len(s) {
			logE.Println("There was an error finding the appropriate element.")
		}

		if s[i] == '[' {
			count++
		} else if s[i] == ']' {
			count--
		}
	}
	return i - 1
}

func translate(s []byte) []Routine {
	const valid = "><+-[].,"
	optimized := make([]Routine, 0)
	switch string(s) {
	// Try to optimize [-] or [+] away
	case "-", "+":
		return append(optimized, Zero)
		// Addition, evaluate [a b] to [0 a+b]
	case "->+<":
		return append(optimized, Plus)
	}

	for i := 0; i < len(s); i++ {
		c := s[i]

		if !strings.Contains(valid, string(c)) {
			continue // Illegal characters should be ignored. This is mandatory to allow our intermediate
			// representation
		}

		switch {
		//// Subtraction, evaluates [x, y] to [x-y, 0]
		//case ss == ">[-<->]<":
		//	optimized = append(optimized, Minus)
		//	i += 7

		// Copy a value
		//case strings.HasPrefix(string(s[i:]), "[->+>+<<]>>[-<<+>>]<<"):
		//	optimized = append(optimized, Copy)
		//	i += 20

		// Exponentiation
		//case strings.HasPrefix(string(s[i:]), ">>+<[->[-<<[->>>+>+<<<<]>>>>[-<<<<+>>>>]<<]>[-<+>]<<]<"):
		//	optimized = append(optimized, Exp)
		//	i += 53

		// Division
		//case strings.HasPrefix(string(s[i:]), "[>[->+>+<<]>[-<<-[>]>>>[<[-<->]<[>]>>[[-]>>+<]>-<]<<]>>>+<<[-<<+>>]<<<]>>>>>[-<<<<<+>>>>>]<<<<<"):
		//	optimized = append(optimized, Divide)
		//	i += 94

		case c == '[':
			loc := beunSearch(s, i)
			if loc >= len(s) {
				logE.Println("FATAL ERROR Trying to optimize loop %s", s[i+1:])
			}
			loop := s[i+1 : loc]
			optimized = append(optimized, Loop{ops: translate(loop)})
			i += len(loop) + 1

		case c == '>', c == '<', c == '-', c == '+':
			// Try to find matching ones afterwards and collapse that
			var count = 1
			for i+count < len(s) && s[i+count] == c {
				count++
			}
			// Count now contains the amount of chars equal to current character
			if count == 1 {
				optimized = append(optimized, toOp(c))
			} else {
				optimized = append(optimized, toOpWithArg(c, count))
				i += int(count) - 1 // Skip over all iterations
				// -1 because the loop adds one
			}

		case c == '.':
			optimized = append(optimized, Print)
		case c == ',':
			optimized = append(optimized, Input)
		}
	}

	// Try to optimize loop

	return optimized
}
