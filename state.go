package main

import (
	"fmt"
	"math"
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
	offset  int8
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
		jumpFwd: make(map[uint]uint, 0),
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
	if statistics {
		s.stats.gt++
	}
}

func (s *State) IndexDec(n uint) { // <
	if s.index >= n {
		s.index -= n
	} else { // We need to reduce the index by the maximum amount.
		s.index = 0
	}

	if statistics {
		s.stats.lt++
	}
}

func (s *State) DataInc(N uint, offset int8) {
	n := uint8(N)
	offsetIndex := int(s.index) + int(offset)
	//s.data[s.index] += n
	if s.data[offsetIndex] < math.MaxUint8-n-1 {
		s.data[offsetIndex] += n
	} else {
		s.data[offsetIndex] = math.MaxUint8 - 1
	}
	if statistics {
		s.stats.plus++
	}
}

func (s *State) DataDec(N uint, offset int8) {
	n := uint8(N)
	offsetIndex := int(s.index) + int(offset)
	//s.data[s.index] -= n
	if s.data[offsetIndex] >= n {
		s.data[offsetIndex] -= n
	} else {
		s.data[offsetIndex] = 0
	}
	if statistics {
		s.stats.minus++
	}
}

func (s *State) StartLoop() {
	if s.data[s.index] != 0 { // Enter loop; save return address on stack
		// Explicitly create copy of instruction counter
		i := s.instr
		s.stack.Push(i) // Push the (pass-by-value) copy.
	} else { // Skip the loop
		s.instr = s.jumpFwd[s.instr]
	}
	if statistics {
		s.stats.startL++
	}
}

func (s *State) EndLoop() {
	//if s.stack.Len() == 0 {
	//	logE.Println("Cannot resolve corresponding bracket, stack is empty")
	//	os.Exit(1)
	//}
	if s.data[s.index] != 0 { // Jump back
		s.instr = s.stack.Get() // Because we add one later
	} else { // End the loop
		s.stack.Pop() // Pop value from stack.
	}
	if statistics {
		s.stats.endL++
	}
}
func (s *State) Print() {
	s.output += string(s.data[s.index])
	fmt.Printf("%c", s.data[s.index])
	if statistics {
		s.stats.dot++
	}
}
func (s *State) Input() {
	var c string
	fmt.Scanf("%c", &c)
	s.data[s.index] = c[0]
	if statistics {
		s.stats.comma++
	}
}

// Advanced BF instructions
func (s *State) Zero() {
	s.data[s.index] = 0
}

func (s *State) Offset(offset int8) {
	s.offset = offset
}

func (s *State) Plus() {
	s.DataInc(uint(s.data[s.index]), 1)
	s.Zero()
}

func (s *State) Minus() {
	s.DataDec(uint(s.data[s.index]+1), 0)
	s.data[s.index+1] = 0
}

func (s *State) Mult() {
	s.DataInc(uint(s.data[s.index]*s.data[s.index+1]), 2)
	s.Zero()
}

func (s *State) Copy() {
	s.DataInc(uint(s.data[s.index]), 1)
}

func (s *State) Exp() {
	s.data[s.index+1] = Pow(s.data[s.index], s.data[s.index+1])
	s.data[s.index+1] = 0
}

func Pow(x uint8, y uint8) (result uint8) {
	var i uint8
	result = 1
	for i = 0; i < y; i++ {
		result *= x
	}
	return
}

func (s *State) Divide() {
	x := s.data[s.index]
	s.data[s.index] = x / s.data[s.index+1]
	s.data[s.index+1] = x % s.data[s.index+1]
}
