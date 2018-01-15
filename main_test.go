package main

import (
	"testing"
)

func TestIncrementIndex(t *testing.T) {
	s := GenState()
	s.IndexInc(1200)
	if s.index != 1200 {
		t.Errorf("Incrementing index failed")
	}
}

func BenchmarkCompleteProgram(b *testing.B) {
	//filename := "test.bf"
	//
	//rawBF, err := ioutil.ReadFile(filename)
	//if err != nil {
	//	logE.Println(err)
	//}
	//
	//bfCode := string(rawBF)

	for n := 0; n < b.N; n++ {
		main()
	}
}

func runHelper(s string) State {
	var state = GenState()
	b := []byte(s)
	bfCode, jumpfwd := translate(&b)
	state.jumpFwd = jumpfwd

	for state.instr < uint(len(bfCode)) {
		bfExecute(&bfCode, &state)
	}
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
