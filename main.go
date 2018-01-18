package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const BUFSIZE = 30000

var (
	logE = log.New(os.Stderr, "[ERRO] ", log.Ldate+log.Ltime+log.Ltime+log.Lshortfile)
	//logW = log.New(os.Stdout, "[WARN] ", log.Ldate+log.Ltime)
	//logI = log.New(os.Stdout, "[INFO] ", log.Ldate+log.Ltime)
)

const statistics = true

func main() {
	startTime := time.Now()
	fmt.Println("Running")

	var filename = flag.String("filename", "test.bf", "Path to file containing BrainYuck code")
	//var doStats = flag.Bool("stats", true, "Disable statistics")
	flag.Parse()
	//statistics = *doStats

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

func NumFormat(N int) string {
	n := int64(N)
	if n < 0 {
		return "-" + NumFormat(-N)
	}
	in := strconv.FormatInt(n, 10)
	out := make([]byte, len(in)+(len(in)-1)/3)

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

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
