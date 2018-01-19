package main

import (
	"fmt"
)

// All state in the program. Executables operate on a state
// This makes all operations unit-testable
type State struct {
	mem    [BUFSIZE]uint8 // Data array used by BF
	ptr    int            // index into the array
	stats  *Stats         // Statistics collected by the operations
	output string         // What all `.` chars produce
}

// Statistics
type Stats struct {
	gt     int
	plus   int
	dot    int
	comma  int
	plusOp int
	minOp  int
	zero   int
	mult   int
	seek   int
	indexL int
	loop   map[string]int
}

func (s *Stats) Sum() int {
	return s.gt +
		s.plus +
		s.dot +
		s.comma +
		s.plusOp +
		s.minOp +
		s.zero +
		s.mult +
		s.seek +
		s.indexL
}

func GenState() State {
	return State{
		//mem: make([]uint8, BUFSIZE),
		//ptr: 0,
		stats: &Stats{
			loop: make(map[string]int),
		},
	}
}

func (s *State) printStats() {
	fmt.Println("Printing statistics")
	fmt.Printf(" gt:     \t%s\n", numFmt(s.stats.gt))
	fmt.Printf(" simplePlus:\t%s\n", numFmt(s.stats.plus))
	fmt.Printf(" dot:    \t%s\n", numFmt(s.stats.dot))
	fmt.Printf(" comma:  \t%s\n", numFmt(s.stats.comma))
	fmt.Printf(" plusOps:\t%s\n", numFmt(s.stats.plusOp))
	fmt.Printf(" minOps: \t%s\n", numFmt(s.stats.minOp))
	fmt.Printf(" zero:   \t%s\n", numFmt(s.stats.zero))
	fmt.Printf(" mult:   \t%s\n", numFmt(s.stats.mult))
	fmt.Printf(" seek:   \t%s\n", numFmt(s.stats.seek))
	fmt.Printf(" loop:  \t%s\n", numFmt(s.stats.indexL))
	fmt.Printf(" Sum number of executions cycles: %s\n", numFmt(s.stats.Sum()))
}

// Normal BF instructions

func (s *State) IndexInc(n int) { // >
	s.ptr += n
	//
	//if s.ptr < 0 {
	//	s.ptr = 0
	//} else if s.ptr >= BUFSIZE {
	//	s.ptr = BUFSIZE - 1
	//}

	if statistics {
		s.stats.gt++
	}
}

func (s *State) DataInc(N int, offset int) {
	n := uint8(N)
	index := int(s.ptr) + offset
	//if index > BUFSIZE || index < 0 {
	//	logE.Printf("Cannot complete this operation, index out of range: index = %d", index)
	//	return
	//}
	s.mem[index] += n

	if statistics {
		s.stats.plus++
	}
}

func (s *State) Print() {
	s.output += string(s.mem[s.ptr])
	if !buffer {
		fmt.Printf("%c", s.mem[s.ptr])
	}
	if statistics {
		s.stats.dot++
	}
}
func (s *State) Input() {
	var c string
	panic("Not yet implemented")
	fmt.Scanf("%c", &c)
	s.mem[s.ptr] = byte(c[0])
	if statistics {
		s.stats.comma++
	}
}

// Advanced BF instructions
func (s *State) Zero() {
	s.mem[s.ptr] = 0
	if statistics {
		s.stats.zero++
	}
}

func (s *State) Plus() {
	s.DataInc(int(s.mem[s.ptr]), 1)
	s.Zero()
	if statistics {
		s.stats.plusOp++
	}
}

// Subtraction, evaluates [x, y] to [x-y, 0]
func (s *State) Minus() {
	s.DataInc(-int(s.mem[s.ptr+1]), 0)
	s.mem[s.ptr+1] = 0

	if statistics {
		s.stats.minOp++
	}
}

func (s *State) Mult() {
	s.DataInc(int(s.mem[s.ptr]*s.mem[s.ptr+1]), 2)
	s.Zero()
	if statistics {
		s.stats.mult++
	}
}

func (s *State) Exp() {
	s.mem[s.ptr+1] = Pow(s.mem[s.ptr], s.mem[s.ptr+1])
	s.mem[s.ptr+1] = 0
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
	x := s.mem[s.ptr]
	s.mem[s.ptr] = x / s.mem[s.ptr+1]
	s.mem[s.ptr+1] = x % s.mem[s.ptr+1]
}

func (s *State) Seek(n int) {
	for s.mem[s.ptr] != 0 {
		s.ptr += n
		//if s.ptr < 0 {
		//	logE.Println("Cannot complete Seek operation")
		//}
	}
	if statistics {
		s.stats.seek++
	}
}

func (s *State) AddAndZero(r []Executable) {
	var i = s.mem[s.ptr]

	for _, o := range r {
		opao, ok := o.(OpWithArgOffset)
		if !ok || opao.op != DataIncArgOffset {
			panic("It's always an op with arg offset in this place.. should be?")
		}
		newIndex := s.ptr + opao.offset
		if newIndex > BUFSIZE || newIndex < 0 {
			continue
		}
		val := opao.arg * int(i) % 256

		s.mem[newIndex] += uint8(val)
	}
	s.mem[s.ptr] = 0
	if statistics {
		s.stats.indexL++
	}
}
