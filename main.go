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

const (
	// Length of the array containing
	BUFSIZE = 30000

	// Whether we should print function call statistics and timing information
	statistics = false

	// Whether output should be printed immediately, or buffered
	buffer = false
)

var (
	logE = log.New(os.Stderr, "[ERRO] ", log.Ldate+log.Ltime+log.Ltime+log.Lshortfile)

	//Do we interpret code, or compile to some language?
	interpret bool
)

func main() {
	startTime := time.Now()
	fmt.Println("Running")

	filename := flag.String("filename", "test.bf", "Path to file containing BrainYuck code")
	compileTo := flag.String("compile", "no", "By default, we interpret the BF code instead of "+
		"compiling to C code. This incurs a penalty hit, compiling is faster.")
	flag.Parse()

	rawBF, err := ioutil.ReadFile(*filename)
	if err != nil {
		logE.Println(err)
		os.Exit(1)
	}
	loadTime := time.Since(startTime)

	var state = GenState()

	// Optimize BF to intermediate representation
	_, ops := translate(rawBF)
	program := Loop{NoLoop, ops}
	initTime := time.Since(startTime) // Timing stuffs

	var compileTime = time.Nanosecond
	switch *compileTo {
	case "c":
		// Compile instructions to C code, execute this C code
		compileTime = compileCAndExec(program, startTime)
	default:
		fmt.Println("This language is not supported at the moment, interpreting")
		fallthrough
	case "no":
		interpret = true
		// Running instructions
		program.execute(&state)

	}

	// Print output
	if buffer && interpret {
		fmt.Printf("\n%s\n", state.output)
	}

	// Timing stuffs
	runTime := time.Since(startTime)
	fmt.Printf("Loading BF code from file:  \t \t\t\t%s\n", loadTime)
	fmt.Printf("Optimizing BF instructions:  \t\t\t\t%s\n", initTime)
	fmt.Printf("Writing to C file, and compiling with %s: \t\t%s\n", optimization, compileTime)
	fmt.Printf("Total duration: \t\t\t\t\t%s\n", runTime)

	// Function call statistics
	if statistics && interpret {
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

			if count == 1 {
				optimized = append(optimized, toOp(c))
				continue
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
		return AddLoop, optimizeGenericAdd(ptr, mem), nil
	}
}

func optimizeGenericAdd(index int, mem map[int]int) (optimized []Executable) {
	// Generic version of AddAndZero loop
	for k, v := range mem {
		if k == 0 {
			// Current cell has to be appended as the last instruction in the loop, otherwise the offsets of the others will be false
			continue
		}
		optimized = append(optimized, OpWithArgOffset{
			op:     DataIncArgOffset,
			arg:    v,
			offset: k,
		})
	}
	optimized = append(optimized,
		// Append zero as last DataInc instruction
		OpWithArg{
			op:  DataIncArg,
			arg: mem[0],
		},
		// Move the pointer the specified amount
		OpWithArg{
			op:  IndexIncArg,
			arg: index,
		},
	)
	return
}

func optimizeAddAndZero(mem map[int]int) (optimized []Executable) {
	// It's an ACTUAL zero ptr loop!!
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
