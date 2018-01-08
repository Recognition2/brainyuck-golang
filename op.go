package main

type Routine interface {
	execute(s *State)
}

type OpWithArgOffset struct {
	op     Op
	arg    uint8
	offset uint8
}

type LoopType uint8

const (
	NoLoop LoopType = iota
	ZeroIndexLoop
)

type Loop struct {
	l   LoopType
	ops []Op
}

func (l Loop) execute(s *State) {
	for _, o := range l.ops {
		o.execute(s)
	}
}

type Op uint8

func (op Op) execute(s *State) {

}

// Define all operations
const (
	// Traditional ops
	DataDec  Op = iota // 0
	DataInc            // 1
	IndexDec           // 2
	IndexInc           // 3

	// Operations that take an argument
	DataDec     // 4
	DataInc     // 5
	IndexDecArg // 6
	IndexIncArg // 7A

	// Operations that use an offset
	DataDecOffset    // 8
	DataIncOffset    // 9
	DataDecArgOffset // 10
	DataIncArgOffset // 11

	Print     // 12
	Input     // 13
	StartLoop // 14 // Special
	EndLoop   // 15

	// Advanced ops
	Zero
	Offset
	Copy
	Plus
	Minus
	Mult
	Exp
	Divide
)

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

func toOpWithArg(c uint8) Op {
	switch c {
	case '<':
		return IndexDecArg
	case '>':
		return IndexIncArg
	case '+':
		return DataIncArg
	case '-':
		return DataDecArg
	default:
		logE.Printf("This is not a valid Op: %d!", c)
		return IndexDec
	}
}

func isOpWithArg(op Op) bool {
	return op >= DataDecArg && op <= IndexIncArg
}
