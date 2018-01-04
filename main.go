package main

import (
	"io/ioutil"
	"log"
	"os"
)

const BUFSIZE = 30000

var (
	logE = log.New(os.Stderr, "[ERRO] ", log.Ldate+log.Ltime+log.Ltime+log.Lshortfile)
	logW = log.New(os.Stdout, "[WARN] ", log.Ldate+log.Ltime)
	logI = log.New(os.Stdout, "[INFO] ", log.Ldate+log.Ltime)
)

func main() {
	filename := "test.bf"

	rawBF, err := ioutil.ReadFile(filename)
	if err != nil {
		logE.Println(err)
	}

	bfCode := string(rawBF)

	var state = GenState()

	// Build map with forward jump addresses
	state.jumpFwd = buildMap(bfCode)

	// Run all instructions
	for state.instr < uint(len(bfCode)) {
		bfExecute(bfCode, &state)
	}
	state.printStats()
	return
}

func buildMap(s string) map[uint]uint {
	//valid := ""
	newMap := make(map[uint]uint, STACKSIZE)
	stack := make(Stack, STACKSIZE)
	for i, c := range s {
		if c == '[' {
			stack.Push(uint(i)) // Store "Begin" address
		} else if c == ']' {
			newMap[stack.Pop()] = uint(i) // map begin address to end address
		}
	}
	return newMap
}

func bfExecute(bf string, state *State) {
	// Endless while
	//fmt.Printf("%c", bf[instr])
	switch bf[state.instr] {
	case '>':
		state.IncrementIndex(bf)
	case '<':
		state.DecrementIndex(bf)
	case '+':
		state.IncrementData()
	case '-':
		state.DecrementData()
	case '.':
		state.Print()
	case ',':
		state.Input()
	case '[':
		state.StartLoop(bf)
	case ']':
		state.EndLoop()
	default:
		state.stats.skipped++
	}
	state.instr++
}
