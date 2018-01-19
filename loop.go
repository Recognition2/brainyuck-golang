package main

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
	default:
		logE.Printf("Is not a valid Loop: op = %d", l.op)
	}
}
