package main

import (
	"testing"
)

func TestStack_Pop(t *testing.T) {
	s := make(Stack, 1)
	s[0] = 120
	if a := s.Pop(); a != 120 {
		t.Errorf("Popping failed. Got %d, want %d", a, 120)
	}
}

func TestStack_Push(t *testing.T) {
	s := make(Stack, 0)
	s.Push(1234)
	if a := s[0]; a != 1234 {
		t.Errorf("Pushing failed. Got %d, want %d", a, 1234)
	}
}

func TestStack_Complete(t *testing.T) {
	s := make(Stack, 0)
	s.Push(120)
	if a := s.Pop(); a != 120 {
		t.Errorf("Cannot reliably store data on the stack, want %d, got %d.", 120, a)
	}
}

func TestIncrementIndex(t *testing.T) {
	s := GenState()
	s.IncrementIndex()
	if s.index != 1 {
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
