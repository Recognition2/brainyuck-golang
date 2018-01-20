package main

import (
	"bytes"
	"fmt"
)

// Operations that take both an argument and an offset, implements Executable
type OpWithArgOffset struct {
	op     Op
	arg    int
	offset int
}

func (o OpWithArgOffset) execute(s *State) {
	switch o.op {
	case DataIncArgOffset:
		s.DataInc(o.arg, o.offset)
	default:
		logE.Printf("Is not a valid Op With Arg And Offset: op = %d", o.op)
	}
}

func (o OpWithArgOffset) toC(b *bytes.Buffer) {
	var c string
	switch o.op {
	case DataIncArgOffset:
		c = fmt.Sprintf("*(ptr + %d) += %d * counter", o.offset, o.arg)
	default:
		logE.Printf("Is not a valid Op With Arg: op = %d", o.op)
		panic("Transformation to C code not implemented")
	}
	b.WriteString(c + ";\n")
}
