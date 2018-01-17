package main

import (
	"errors"
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

	// Optimize BF to intermediate representation
	_, ops := translate(rawBF)
	program := Loop{NoLoop, ops}

	// Optimize loops even harder
	//loop = optimizeLoops(bfCode, jumpfwd)
	fmt.Println("Done optimizing, running...")
	initTime := time.Since(startTime)

	// Run all instructions
	program.execute(&state)

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

//func optimizeLoops(loop []Op) []Op {
//	var i uint
//	for i = 0; i < uint(len(loop)); i++ {
//		op := loop[i]
//		if op == StartLoop {
//			end := skipLoopMap[i]
//			content := loop[i:end]
//			whatLoop, whatHappensMap := analyseLoop(content)
//			var newOps = make([]Op, 0)
//			switch whatLoop {
//			case NoLoop:
//				continue
//			case ZeroIndexLoop:
//				newOps = optimizeZeroIndex(whatHappensMap)
//			default:
//			}
//			//logE.Printf("Current operation is GOOD: %v\n", loop[end] == EndLoop)
//			//logE.Printf("Correct one is -1: %v, +1: %v", loop[end-1] == EndLoop, loop[end+1] == EndLoop)
//			//panic("aaaaa")
//
//			// If the loop was optimized, replace the old code with the new code.
//			loop = append(append(loop[:i], newOps...), loop[end+1:]...)
//		}
//	}
//	return loop
//}
//
//func analyseLoop(loop []Op) (Op, map[int]int) {
//
//	for i := 0; i < len(loop); i++ {
//		op := loop[i]
//		switch op {
//		case Print, Input:
//			return NoLoop, nil
//		case StartLoop, EndLoop:
//			return NoLoop, nil
//		case IndexDec:
//			index--
//		case IndexDecArg:
//			index -= int(loop[i+1])
//			i++
//		case IndexInc:
//			index++
//		case IndexIncArg:
//			index += int(loop[i+1])
//			i++
//		case DataDec:
//			change[index]--
//		case DataDecArg:
//			change[index] -= int(loop[i+1])
//		case DataInc:
//			change[index]++
//		case DataIncArg:
//			change[index] += int(loop[i+1])
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

func translate(s []byte) (Op, []Routine) {
	const valid = "><+-[].,"
	optimized := make([]Routine, 0)
	switch string(s) {
	// Try to optimize [-] or [+] away
	case "-", "+":
		return NoLoop, append(optimized, Zero)

	// Addition, evaluate [a b] to [0 a+b]
	case "->+<":
		return NoLoop, append(optimized, Plus)
	}

	// Optimize ZeroIndexLoop
	ops, err := tryOptimizeZeroIndexLoop(s)
	if err == nil {
		return ZeroIndexLoop, ops
	}

	// Default loop handling
	for i := 0; i < len(s); i++ {
		c := s[i]

		if !strings.Contains(valid, string(c)) {
			continue // Illegal characters should be ignored. This is mandatory to allow our intermediate
			// representation
		}

		switch c {
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

		case '[':
			loc := beunSearch(s, i)
			if loc >= len(s) {
				logE.Println("FATAL ERROR Trying to optimize loop %s", s[i+1:])
			}
			loop := s[i+1 : loc]
			op, instr := translate(loop)
			optimized = append(optimized, Loop{op: op, loop: instr})
			i += len(loop) + 1

		case '>', '<', '-', '+':
			// Try to find matching ones afterwards and collapse that
			var count = 1
			for i+count < len(s) && s[i+count] == c {
				count++
			}
			// Count now contains the amount of chars equal to current character
			//if count == 1 {
			//	optimized = append(optimized, toOp(c))
			//} else {
			//
			if len(s) == count && (c == '>' || c == '<') { // Seek operation
				if c == '<' {
					count = -count
				}
				return NoLoop, append(make([]Routine, 0), OpWithArg{op: Seek, arg: count})
			}

			optimized = append(optimized, toOpWithArg(c, count))
			i += int(count) - 1 // Skip over all iterations
			// -1 because the loop adds one
			//}

		case '.':
			optimized = append(optimized, Print)
		case ',':
			optimized = append(optimized, Input)
		}
	}

	return DefaultLoop, optimized
}
func tryOptimizeZeroIndexLoop(s []byte) ([]Routine, error) {
	index := 0
	change := make(map[int]int)

	for _, c := range s {
		switch c {
		case '.', ',', '[', ']':
			return nil, errors.New("This is not a zero index loop")
		case '+':
			change[index] += 1
		case '-':
			change[index] -= 1
		case '>':
			index++
		case '<':
			index--
		default:
			continue
		}
	}

	if !(index == 0 && change[0] == -1) {
		return nil, errors.New("Not zero index")
	}

	// It's an ACTUAL zero index loop!!
	optimized := make([]Routine, 0)
	for k, v := range change {
		if k == 0 {
			continue
		}
		optimized = append(optimized, OpWithArgOffset{
			op:     DataIncArgOffset,
			arg:    v,
			offset: k,
		})
	}
	return optimized, nil
}
