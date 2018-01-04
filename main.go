package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const BUFSIZE = 30000

var (
	logE = log.New(os.Stderr, "[ERRO] ", log.Ldate+log.Ltime+log.Ltime+log.Lshortfile)
	logW = log.New(os.Stdout, "[WARN] ", log.Ldate+log.Ltime)
	logI = log.New(os.Stdout, "[INFO] ", log.Ldate+log.Ltime)
)

func main() {
	fmt.Println("Running")

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

	fmt.Println("Done optimizing, running...")
	// Run all instructions
	for state.instr < uint(len(bfCode)) {
		bfExecute(bfCode, &state)
	}

	fmt.Printf("\n%s\n", state.output)
	state.printStats()
	return
}

func optimize(s []byte) ([]byte, map[uint]uint) {
	valid := "><+-[].,"
	optimizedCode := bytes.Buffer{}

	skipLoopMap := make(map[uint]uint, STACKSIZE) // Build map for skipping loops
	stack := make(Stack, 0)

	newIndex := 0
	for i := 0; i < len(s); i++ {
		c := s[i]

		if !strings.Contains(valid, string(c)) {
			continue // Illegal character
		}

		switch c {
		case '[':
			// Try to optimize `[-]` and `[+]` away
			if s[i+2] == ']' && (s[i+1] == '-' || s[i+1] == '+') {
				optimizedCode.WriteByte('0')
				newIndex++
				i += 2
				continue
			}
			stack.Push(uint(newIndex)) // Store "Begin" address
			optimizedCode.WriteByte(c)
		case ']':
			skipLoopMap[stack.Pop()] = uint(newIndex) // map begin address to end address
			optimizedCode.WriteByte(c)
		case '>', '<', '-', '+':
			// Try to find matching ones afterwards and collapse that
			count := 0
			for s[i+count] == c && count < 255 {
				count++
			} // Count now contains the amount of chars equal to current character
			optimizedCode.WriteByte(c)
			optimizedCode.WriteByte(byte(count)) // This is an illegal char according to BF parser
			newIndex++                           // Skip the "amount" byte
			i += count - 1
		default:
			optimizedCode.WriteByte(c)
		}
		newIndex++
	}
	return optimizedCode.Bytes(), skipLoopMap
}

func executeWithCount(c byte, n uint, state *State) {
	switch c {
	case '>':
		state.IncrementIndex(n)
	case '<':
		state.DecrementIndex(n)
	case '+':
		state.IncrementData(uint8(n))
	case '-':
		state.DecrementData(uint8(n))
	default:
		logE.Println("Cannot execute irrelevant byte")
	}
}

func bfExecute(bf []byte, state *State) {
	switch c := bf[state.instr]; c {
	case '>', '<', '+', '-':
		n := bf[state.instr+1]
		executeWithCount(c, uint(n), state)
		//fmt.Printf("%c%d ", c, bf[state.instr+1])
		state.instr += 2
		return
	case '.':
		state.Print()
	case ',':
		state.Input()
	case '[':
		state.StartLoop(bf)
	case ']':
		state.EndLoop()
	case '0':
		state.Zero()
	default:
		logE.Printf("This cannot happen: c = %d", c)
	}
	//fmt.Printf("%c ", bf[state.instr])
	state.instr++
}
