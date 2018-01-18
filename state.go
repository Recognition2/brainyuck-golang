package main

import (
	"fmt"
)

const STACKSIZE = 500

type State struct {
	data  []uint8
	index int
	stack Stack
	//instr  uint
	stats  Stats
	output string
	offset int
}

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

func (s Stats) Total() int {
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
		// Data buffer
		data:  make([]uint8, BUFSIZE),
		index: 0,
		// Stack to correctly implement loops
		stack: make(Stack, 0),
		//instr: 0,
		stats: Stats{loop: make(map[string]int)},
	}
}

func (s *State) printStats() {
	fmt.Println("Printing statistics")
	fmt.Printf(" gt:     \t%s\n", numFormat(s.stats.gt))
	//fmt.Printf(" lt:     \t%s\n", numFormat(s.stats.lt))
	fmt.Printf(" simplePlus:\t%s\n", numFormat(s.stats.plus))
	fmt.Printf(" dot:    \t%s\n", numFormat(s.stats.dot))
	fmt.Printf(" comma:  \t%s\n", numFormat(s.stats.comma))
	fmt.Printf(" plusOps:\t%s\n", numFormat(s.stats.plusOp))
	fmt.Printf(" minOps: \t%s\n", numFormat(s.stats.minOp))
	fmt.Printf(" zero:   \t%s\n", numFormat(s.stats.zero))
	fmt.Printf(" mult:   \t%s\n", numFormat(s.stats.mult))
	fmt.Printf(" seek:   \t%s\n", numFormat(s.stats.seek))
	fmt.Printf(" index:  \t%s\n", numFormat(s.stats.indexL))
	fmt.Printf(" Total number of executions cycles: %s\n", numFormat(s.stats.Total()))
}

func (s *State) PrintState() {
	fmt.Print("Logging s: \n Data: ")

	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", s.data[i])
	}

	fmt.Printf("\n Current index: %d\n", s.index)
	//fmt.Printf(" Current instr: %d\n", s.instr)
}

// Normal BF instructions

func (s *State) IndexInc(n int) { // >
	s.index += n

	//if n < 0 {
	//	n = 0
	//} else if n >= BUFSIZE {
	//	n = BUFSIZE - 1
	//}

	if statistics {
		s.stats.gt++
	}
}

func (s *State) DataInc(N int, offset int) {
	n := uint8(N)
	//offsetIndex := int(s.index) + offset
	//s.data[s.index] += n
	offset = int(s.index) + offset
	s.data[offset] += n
	//if s.data[offsetIndex] < math.MaxUint8-n-1 {
	//	s.data[offsetIndex] += n
	//} else {
	//	s.data[offsetIndex] = math.MaxUint8 - 1
	//}
	if statistics {
		s.stats.plus++
	}
}

func (s *State) Print() {
	s.output += string(s.data[s.index])
	//fmt.Printf("%c", s.data[s.index])
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
	if statistics {
		s.stats.zero++
	}
}

func (s *State) Plus() {
	s.DataInc(int(s.data[s.index]), 1)
	s.Zero()
	if statistics {
		s.stats.plusOp++
	}
}

func (s *State) Minus() {
	s.DataInc(-int(s.data[s.index+1]), 0)
	s.data[s.index+1] = 0
	if statistics {
		s.stats.minOp++
	}
}

func (s *State) Mult() {
	s.DataInc(int(s.data[s.index]*s.data[s.index+1]), 2)
	s.Zero()
	if statistics {
		s.stats.mult++
	}
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

func (s *State) Seek(n int) {
	for s.data[s.index] != 0 {
		s.index += n
		//if s.index < 0 {
		//	logE.Println("Cannot complete Seek operation")
		//}
	}
	if statistics {
		s.stats.seek++
	}
}

func (s *State) ZeroIndexLoop(r []Routine) {
	var i = s.data[s.index]

	for _, o := range r {
		opao, ok := o.(OpWithArgOffset)
		if !ok || opao.op != DataIncArgOffset {
			panic("It's always an op with arg offset in this place.. should be?")
		}
		newIndex := s.index + opao.offset
		if newIndex > BUFSIZE || newIndex < 0 {
			continue
		}
		val := opao.arg * int(i) % 256

		s.data[newIndex] += uint8(val)
	}
	s.data[s.index] = 0
	if statistics {
		s.stats.indexL++
	}
}
