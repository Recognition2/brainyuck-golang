package main

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
