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

func GenState() State {
	return State{
		// Data buffer
		data:  make([]uint8, BUFSIZE),
		index: 0,
		// Stack to correctly implement loops
		stack:   make(Stack, STACKSIZE),
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
}
func (s *State) PrintState() {
	fmt.Print("Logging s: \n Data: ")

	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", s.data[i])
	}

	fmt.Printf("\n Current index: %d\n", s.index)
	fmt.Printf(" Current instr: %d\n", s.instr)
}

func (s *State) IncrementIndex(bf string) { // >
	if s.index < BUFSIZE {
		s.index++
	}
	s.stats.gt++
}

func (s *State) DecrementIndex(bf string) { // <
	if s.index > 0 {
		s.index--
	}
	s.stats.lt++
}

func (s *State) IncrementData() {
	if s.data[s.index] < math.MaxUint8 {
		s.data[s.index]++
	}
	s.stats.plus++
}

func (s *State) DecrementData() {
	if s.data[s.index] > 0 {
		s.data[s.index]--
	}
	s.stats.minus++
}

func (s *State) StartLoop(bf string) {
	// Normal (non-opt) behaviour

	if s.data[s.index] != 0 { // Enter loop; save return address on stack
		s.stack.Push(s.instr)
	} else { // Skip the loop
		s.instr = s.jumpFwd[s.instr]
	}
	s.stats.startL++
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
	fmt.Printf("%c", s.data[s.index])
	s.stats.dot++
}
func (s *State) Input() {
	var c string
	fmt.Scanf("%c", &c)
	s.data[s.index] = c[0]
	s.stats.comma++
}
