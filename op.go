package main

// Define all operations
const (
	// Traditional ops
	DataDec  Op = iota // 0
	DataInc            // 1
	IndexDec           // 2
	IndexInc           // 3

	Print // 12
	Input // 13
	//StartLoop // 14 // Special
	//EndLoop   // 15

	// Advanced ops
	Zero
	Plus
	Minus
	Mult
	Exp
	Divide
	Seek

	// Ops with arguments
	DataIncArg // ArgOp = iota
	IndexIncArg

	// Ops with arguments and offsets
	DataIncArgOffset // ArgOffsetOp = iota

	// Loops
	NoLoop
	DefaultLoop
	ZeroIndexLoop
)

type Routine interface {
	execute(s *State)
}

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
		logE.Printf("This cannot happen: op = %d", o.op)
	}
}

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
		logE.Printf("This cannot happen: op = %d", o.op)
	}
}

type Loop struct {
	op   Op
	loop []Routine
	//counter *int
}

func (l Loop) execute(s *State) {
	switch l.op {
	case DefaultLoop:
		for s.data[s.index] != 0 {
			for _, o := range l.loop {
				o.execute(s)
			}
			//if statistics {
			//	*l.counter++
			//}
		}
	case NoLoop:
		for _, o := range l.loop {
			o.execute(s)
		}
	case ZeroIndexLoop:
		s.ZeroIndexLoop(l.loop)
	}
}

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
		logE.Printf("This cannot happen: op = %d", op)
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

func toOpWithArg(c uint8, count int) Routine {
	switch c {
	case '>':
		return OpWithArg{op: IndexIncArg, arg: count}
	case '<':
		a := OpWithArg{op: IndexIncArg, arg: -count}
		return a
	case '+':
		return OpWithArg{op: DataIncArg, arg: count}
	case '-':
		return OpWithArg{op: DataIncArg, arg: -count}
	default:
		logE.Printf("This is not a valid Routine: %d!", c)
		return IndexInc
	}
}
