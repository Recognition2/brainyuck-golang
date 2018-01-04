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
