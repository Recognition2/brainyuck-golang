package main

import "bytes"

/// This file contains all Operation-related stuffs.

// Standard Op type, does not take an argument.
type Op uint

// All BF IML operations.
const (
	// Traditional ops, known in BF syntax
	DataDec  Op = iota // 0 -
	DataInc            // 1 +
	IndexDec           // 2 <
	IndexInc           // 3 >
	Print              // 4 .
	Input              // 5 ,

	// Advanced ops, recognized from traditional BF operations
	Zero   // 6 [-], [+]
	Plus   // 7 [->+<]
	Minus  // 8
	Mult   // 9
	Exp    // 10
	Divide // 11
	Seek   // 12

	// Ops with arguments
	DataIncArg
	IndexIncArg

	// Ops with arguments and offsets
	DataIncArgOffset // Used only for AddAndZero

	// Loopable operations
	NoLoop      // Just a sequence of operations
	DefaultLoop // Any loop
	AddAndZero  // Not actually a loop, but does carry an array of operations with him.
	AddLoop
)

func (op Op) execute(s *State) {
	switch op {
	case IndexDec:
		s.IndexInc(-1)
	case IndexInc:
		s.IndexInc(1)
	case DataDec:
		s.DataInc(-1, 0)
	case DataInc:
		s.DataInc(1, 0)
	case Print:
		s.Print()
	case Input:
		s.Input()

		// Special operations
	case Zero:
		s.Zero()
	case Plus:
		s.Plus()
	case Minus:
		s.Minus()
	case Mult:
		s.Mult()
	case Exp:
		s.Exp()
	case Divide:
		s.Divide()

	default:
		logE.Printf("Op is not an op: %d", op)
	}
}

func (op Op) toC(b *bytes.Buffer) {
	var c string
	switch op {
	case IndexDec:
		c = "--ptr"
	case IndexInc:
		c = "++ptr"
	case DataDec:
		c = "-- *ptr"
	case DataInc:
		c = "++ *ptr"
	case Print:
		c = `putchar(*ptr)`
	case Input:
		c = `*ptr = getchar()`

		// Special operations
	case Zero:
		c = "*ptr = 0"
	case Plus:
		c = "*(ptr+1) += *ptr;\n" +
			"*ptr = 0"
	case Minus:
		c = "*ptr -= *(ptr+1);\n" +
			"*(ptr+1) = 0"
	case Mult:
		panic("Transformation to C code not implemented")
	case Exp:
		panic("Transformation to C code not implemented")
	case Divide:
		panic("Transformation to C code not implemented")

	default:
		logE.Printf("Op is not an op: %d", op)
		panic("Transformation to C code not implemented")
	}
	b.WriteString(c + ";\n") // Bc otherwise I'm gonna forget those pesky `;` anyway heheh
}

func toOp(c uint8) Op {
	switch c {
	case '<':
		return IndexDec
	case '>':
		return IndexInc
	case '+':
		return DataInc
	case '-':
		return DataDec
	case '.':
		return Print
	case ',':
		return Input
	default:
		logE.Printf("This is not a valid Op: %d!", c)
		return IndexDec
	}
}
