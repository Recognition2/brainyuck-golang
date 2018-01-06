package main

import (
	"fmt"
	"math"
	"os"
)

const STACKSIZE = 500

type State struct {
	data    []uint8
	index   uint
	stack   Stack
	jumpFwd map[uint]uint
	instr   uint
	stats   Stats
	output  string
}

type Stats struct {
	lt      int
	gt      int
	plus    int
	minus   int
	dot     int
	comma   int
	startL  int
	endL    int
	skipped int
}

func (s Stats) Total() int {
	return s.lt +
		s.gt +
		s.plus +
		s.minus +
		s.dot +
		s.comma +
		s.startL +
		s.endL +
		s.skipped
}

func GenState() State {
	return State{
		// Data buffer
		data:  make([]uint8, BUFSIZE),
		index: 0,
		// Stack to correctly implement loops
		stack:   make(Stack, 0),
		jumpFwd: make(map[uint]uint, STACKSIZE),
		instr:   0,
	}
}

func (s *State) printStats() {
	fmt.Println("Printing statistics")
	fmt.Printf(" gt:     \t%d\n", s.stats.gt)
	fmt.Printf(" lt:     \t%d\n", s.stats.lt)
	fmt.Printf(" plus:   \t%d\n", s.stats.plus)
	fmt.Printf(" minus:  \t%d\n", s.stats.minus)
	fmt.Printf(" dot:    \t%d\n", s.stats.dot)
	fmt.Printf(" comma:  \t%d\n", s.stats.comma)
	fmt.Printf(" start L:\t%d\n", s.stats.startL)
	fmt.Printf(" End L:  \t%d\n", s.stats.endL)
	fmt.Printf(" Skipped:\t%d\n", s.stats.skipped)
	fmt.Printf(" Total number of executions cycles: %d\n", s.stats.Total())
}

func (s *State) PrintState() {
	fmt.Print("Logging s: \n Data: ")

	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", s.data[i])
	}

	fmt.Printf("\n Current index: %d\n", s.index)
	fmt.Printf(" Current instr: %d\n", s.instr)
}

// Normal BF instructions

func (s *State) IndexInc(n uint) { // >
	if s.index < BUFSIZE-n {
		s.index += n
	} else {
		s.index = BUFSIZE - 1
	}
	s.stats.gt++
}

func (s *State) IndexDec(n uint) { // <
	if s.index >= n {
		s.index -= n
	} else { // We need to reduce the index by the maximum amount.
		s.index = 0
	}

	s.stats.lt++
}

func (s *State) DataInc(N uint) {
	n := uint8(N)
	if s.data[s.index] < math.MaxUint8-n-1 {
		s.data[s.index] += n
	} else {
		s.data[s.index] = math.MaxUint8 - 1
	}
	s.stats.plus++
}

func (s *State) DataDec(N uint) {
	n := uint8(N)
	if s.data[s.index] >= n {
		s.data[s.index] -= n
	} else {
		s.data[s.index] = 0
	}
	s.stats.minus++
}

func (s *State) StartLoop(bf []Op) {
	if s.data[s.index] != 0 { // Enter loop; save return address on stack
		s.stack.Push(s.instr)
	} else { // Skip the loop
		s.instr = s.jumpFwd[s.instr]
	}
}

func (s *State) EndLoop() {
	if len(s.stack) == 0 {
		logE.Println("Cannot resolve corresponding bracket, stack is empty")
		os.Exit(1)
	}
	if s.data[s.index] != 0 { // Jump back
		s.instr = s.stack.Get() // Because we add one later
	} else { // End the loop
		s.stack.Pop() // Pop value from stack.
	}
	s.stats.endL++
}
func (s *State) Print() {
	s.output += string(s.data[s.index])
	//fmt.Printf("%c", s.data[s.index])
	s.stats.dot++
}
func (s *State) Input() {
	var c string
	fmt.Scanf("%c", &c)
	s.data[s.index] = c[0]
	s.stats.comma++
}

// Advanced BF instructions
func (s *State) Zero() {
	s.data[s.index] = 0
}

func (s *State) Plus() {
	s.data[s.index+1] += s.data[s.index]
	s.data[s.index] = 0
}

func (s *State) Minus() {
	s.DataDec(uint(s.data[s.index] + 1))
	s.data[s.index+1] = 0
}

func (s *State) Mult() {
	s.data[s.index+2] = s.data[s.index] * s.data[s.index+1]
	s.data[s.index] = 0
}
