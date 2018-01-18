package main

import (
	"testing"
)

func TestIncrementIndex(t *testing.T) {
	s := GenState()
	s.IndexInc(1200)
	if s.ptr != 1200 {
		t.Errorf("Incrementing ptr failed")
	}
}

func BenchmarkCompleteProgram(b *testing.B) {
	for n := 0; n < b.N; n++ {
		main()
	}
}

func runHelper(s string) State {
	var state = GenState()
	b := []byte(s)
	_, ops := translate(b)
	program := Loop{NoLoop, ops}
	program.execute(&state)
	return state
}

func TestIncrement(t *testing.T) {
	prog := "++++."
	state := runHelper(prog)
	if state.output != string(4) {
		t.Errorf("Cannot Increment reliably")
	}
}

func TestLoops(t *testing.T) {
	prog := "+++[>++<-]>."
	state := runHelper(prog)
	if state.output != string(6) {
		t.Errorf("Error in the way loops are created, wanted %d found %d", 6, int(state.output[0]))
	}
}

func BenchmarkSmallProgram(b *testing.B) {
	prog := "+++[>++<-]>."
	for n := 0; n < b.N; n++ {
		runHelper(prog)
	}
}
