package main

import (
	"bytes"
	"fmt"
)

// Operations that take an argument. Implements Executable
type OpWithArg struct {
	op  Op
	arg int
}

func (o OpWithArg) execute(s *State) {
	switch o.op {
	case Seek:
		s.Seek(o.arg)
	case IndexIncArg:
		s.IndexInc(o.arg)
	case DataIncArg:
		s.DataInc(o.arg, 0)
	default:
		logE.Printf("Is not a valid Op With Arg: op = %d", o.op)
	}
}

func (o OpWithArg) toC(b *bytes.Buffer) {
	var c string
	switch o.op {
	case Seek:
		c = fmt.Sprintf("while (*ptr) {\n"+
			"ptr += %d;\n}\n", o.arg)
	case IndexIncArg:
		c = fmt.Sprintf("ptr += %d;\n", o.arg)
	case DataIncArg:
		c = fmt.Sprintf("*ptr += %d;\n", o.arg)
	default:
		logE.Printf("Is not a valid Op With Arg: op = %d;\n", o.op)
		panic("Transformation to C code not implemented")
	}
	b.WriteString(c)
}

// toOpWithArg transforms a BF instruction that happens `count` times into an appropriate IML operation
func toOpWithArg(c uint8, count int) OpWithArg {
	switch c {
	case '>':
		return OpWithArg{op: IndexIncArg, arg: count}
	case '<':
		return OpWithArg{op: IndexIncArg, arg: -count}
	case '+':
		return OpWithArg{op: DataIncArg, arg: count}
	case '-':
		return OpWithArg{op: DataIncArg, arg: -count}
	default:
		logE.Printf("This is not a valid Executable: %d!", c)
		return OpWithArg{}
	}
}
