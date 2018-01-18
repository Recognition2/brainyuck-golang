package main

/// This file contains all Operation-related stuffs.

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

// Define interface implemented by all operations.
// An executable operates on a State.
type Executable interface {
	execute(s *State)
}

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

// Sequence of operations that, in general, executes until mem[ptr] reaches zero. Implements Executable
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
		s.GenericAdd(l.loop)
	}
}

// Standard Op type, does not take an argument.
type Op uint

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

func toOp(c uint8) Op {
	switch c {
	case '<':
		return IndexDec
	case '>':
		return IndexInc
	case '+':
		return DataDec
	case '-':
		return DataInc
	case '.':
		return Print
	case ',':
		return Input
	default:
		logE.Printf("This is not a valid Op: %d!", c)
		return IndexDec
	}
}

// toOpWithArg transforms a BF instruction that happens `count` times into an appropriate IML operation
func toOpWithArg(c uint8, count int) Executable {
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
		return IndexInc
	}
}
