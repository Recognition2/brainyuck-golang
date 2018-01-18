package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Length of the array containing
const BUFSIZE = 30000

var (
	logE = log.New(os.Stderr, "[ERRO] ", log.Ldate+log.Ltime+log.Ltime+log.Lshortfile)
	//logW = log.New(os.Stdout, "[WARN] ", log.Ldate+log.Ltime)
	//logI = log.New(os.Stdout, "[INFO] ", log.Ldate+log.Ltime)
)

// Print function call statistics and timing information
const statistics = false
const buffer = false

func main() {
	startTime := time.Now()
	fmt.Println("Running")

	var filename = flag.String("filename", "test.bf", "Path to file containing BrainYuck code")
	//buffer = *flag.Bool("buffer", false, "Whether the output should be buffered until the end")
	flag.Parse()

	rawBF, err := ioutil.ReadFile(*filename)
	if err != nil {
		logE.Println(err)
	}

	var state = GenState()

	// Optimize BF to intermediate representation
	_, ops := translate(rawBF)
	program := Loop{NoLoop, ops}

	runtime.GC()

	// Optimize loops even harder
	fmt.Println("Done optimizing, running...")
	initTime := time.Since(startTime)

	// Run all instructions
	program.execute(&state)

	// Print output
	if buffer {
		fmt.Printf("\n%s\n", state.output)
	}

	// Timing stuffs
	runTime := time.Since(startTime)
	fmt.Printf("Optimizing took %s.\nTotal took %s\n", initTime, runTime)

	// Function call statistics
	if statistics {
		state.printStats()
	}
	return
}

// numFmt formats very long integers as a string
func numFmt(N int) string {
	n := int64(N)
	if n < 0 {
		return "-" + numFmt(-N)
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

// beunSearch finds the corresponding `]` to an `[` at ptr i to the byte array s
// It's called beunSearch because it's an inefficient way of searching
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

// translate transforms a byte array of BF instructions to an IML representation.
// Loops are, if possible, optimized away, or translated to the other representation
func translate(s []byte) (Op, []Executable) {
	optimized := make([]Executable, 0) // Create slice of Executables that the
	ss := string(s)
	switch ss {
	// Try to optimize [-] or [+] away
	case "-", "+":
		return NoLoop, append(optimized, Zero)

	// Addition, evaluate [a b] to [0 a+b]
	case "->+<":
		return NoLoop, append(optimized, Plus)
	}

	// Optimize AddAndZero loops, a special type of loop that loops `n` times, with `n` the mem at the current pointer.
	loopType, ops, err := tryOptimizeAddLoop(s)
	if err == nil {
		return loopType, ops
	}

	// Default loop handling, iterate over contents
	for i := 0; i < len(s); i++ {
		c := s[i]
		// Subtraction, evaluates [x, y] to [x-y, 0]
		if strings.HasPrefix(string(s[i:]), ">[-<->]<") {
			optimized = append(optimized, Minus)
			i += 7
			continue
		} else if strings.HasPrefix(string(s[i]), ">>+<[->[-<<[->>>+>+<<<<]>>>>[-<<<<+>>>>]<<]>[-<+>]<<]<") {
			optimized = append(optimized, Exp)
			i += 53
			continue
		} else if strings.HasPrefix(string(s[i]), "[>[->+>+<<]>[-<<-[>]>>>[<[-<->]<[>]>>[[-]>>+<]>-<]<<]>>>+<<[-<<+>>]<<<]>>>>>[-<<<<<+>>>>>]<<<<<") {
			optimized = append(optimized, Divide)
			i += 94
			continue
		}

		switch c {
		case '[': // Start of loop symbol
			loc := beunSearch(s, i)
			if loc >= len(s) {
				logE.Printf("FATAL ERROR Trying to optimize loop %s", s[i+1:])
			}
			loop := s[i+1 : loc]         // Instructions inside of the nested loop
			op, instr := translate(loop) // Recursive call on the inner loop
			optimized = append(optimized, Loop{op: op, loop: instr})
			i += len(loop) + 1 // +1 to also go beyond the matching `]`

		case '>', '<', '-', '+':
			// Try to find matching ones afterwards and collapse that into 1 operation
			var count = 1
			for i+count < len(s) && s[i+count] == c {
				count++
			}
			// Count now contains the amount of chars equal to current character
			if len(s) == count && (c == '>' || c == '<') { // Entire loop consists of < or >
				if c == '<' {
					count = -count
				}
				return NoLoop, append(make([]Executable, 0), OpWithArg{op: Seek, arg: count})
			}

			optimized = append(optimized, toOpWithArg(c, count))
			i += int(count) - 1 // Skip over all iterations
			// -1 because the loop adds one

		case '.':
			optimized = append(optimized, Print)

		case ',':
			optimized = append(optimized, Input)

		default:
			// Illegal characters should be ignored. This is mandatory to follow any BF spec
			continue
		}
	}
	return DefaultLoop, optimized
}

// tryOptimizeAddLoop tries to optimize a loop into a AddAndZeroLoop, if it is one.
// Otherwise, it tries to find a generic Add loop.
// If not, it returns an error
func tryOptimizeAddLoop(s []byte) (Op, []Executable, error) {
	ptr := 0
	// Changes that happen inside of the loop
	mem := make(map[int]int)

	for _, c := range s {
		switch c {
		case '.', ',', '[', ']':
			return Input, nil, errors.New("this is not a zero ptr loop")
		case '+':
			mem[ptr] += 1
		case '-':
			mem[ptr] -= 1
		case '>':
			ptr++
		case '<':
			ptr--
		default:
			continue
		}
	}

	if ptr == 0 && mem[0] == -1 {
		return AddAndZero, optimizeAddAndZero(mem), nil
	} else {
		// We can, at the moment, not optimize this further.
		return Input, nil, errors.New("cannot optimize generic loops yet.")
	}
}

func optimizeAddAndZero(mem map[int]int) []Executable {
	// It's an ACTUAL zero ptr loop!!
	optimized := make([]Executable, 0)
	for k, v := range mem {
		if k == 0 {
			// The current cell is zeroed anyway, property of this type of loop
			continue
		}
		optimized = append(optimized, OpWithArgOffset{
			op:     DataIncArgOffset,
			arg:    v,
			offset: k,
		})
	}
	return optimized
}
