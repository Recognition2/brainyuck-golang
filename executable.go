package main

// Define interface implemented by all operations.
// An executable operates on a State.
type Executable interface {
	execute(s *State)
}
