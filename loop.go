package main

import (
	"bytes"
)

// Sequence of operations that, in general, executes until mem[ptr] reaches zero.
// Implements Executable
// Does not need to be loopable.
type Loop struct {
	op   Op
	loop []Executable
}

func (l Loop) execute(s *State) {
	switch l.op {
	case DefaultLoop:
		for s.mem[s.ptr] != 0 {
			for _, o := range l.loop {
				o.execute(s)
			}
		}
	case NoLoop:
		for _, o := range l.loop {
			o.execute(s)
		}
	case AddAndZero:
		s.AddAndZero(l.loop)
	case AddLoop:
		for s.mem[s.ptr] != 0 {
			s.GenericAdd(l.loop)
		}
	default:
		logE.Printf("Is not a valid Loop: op = %d", l.op)
	}
}

func (l Loop) toC(b *bytes.Buffer) {
	switch l.op {
	case DefaultLoop:
		b.WriteString("while (*ptr) { \n")
		for _, o := range l.loop {
			o.toC(b)
		}
		b.WriteString("\n}\n")
	case NoLoop:
		for _, o := range l.loop {
			o.toC(b)
		}
	case AddAndZero:
		b.WriteString("counter = (int) *ptr;\n" +
			"*ptr = 0;\n")
		for _, o := range l.loop {
			o.toC(b)
		}
	case AddLoop:
		b.WriteString("while (*ptr) { \n" +
			"counter = 1;\n")
		for _, o := range l.loop {
			o.toC(b)
		}
		b.WriteString("\n}\n")
	default:
		panic("Transformation to C code not yet implemented")
	}
}
