package main

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
