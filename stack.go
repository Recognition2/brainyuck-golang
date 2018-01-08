package main

type Stack []uint

func (s *Stack) Push(v uint) {
	*s = append(*s, v)
}

func (s *Stack) Pop() uint {
	l := len(*s)
	res := (*s)[l-1]
	*s = (*s)[:l-1]
	return res
}

func (s *Stack) Get() uint {
	return (*s)[len(*s)-1]
}

//// Stack is a basic LIFO stack that resizes as needed.
//type Stack struct {
//	nodes []*uint
//	count int
//}
//
//// Push adds a node to the stack.
//func (s *Stack) Push(n *uint) {
//	s.nodes = append(s.nodes[:s.count], n)
//	s.count++
//}
//
//// Pop removes and returns a node from the stack in last to first order.
//func (s *Stack) Pop() *uint {
//	if s.count == 0 {
//		return nil
//	}
//	s.count--
//	return s.nodes[s.count]
//}
//
//// Last returns the last element of the slice
//func (s *Stack) Last() *uint {
//	if s.count == 0 {
//		return nil
//	}
//	return s.nodes[s.count-1]
//}
//
//func NewStack() *Stack {
//	return &Stack{}
//}
