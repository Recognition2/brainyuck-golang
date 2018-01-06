package main

type Op uint8

// Define all operations
const (
	// Traditional ops
	DataDec   Op = iota // 0
	DataInc             // 1
	IndexDec            // 2
	IndexInc            // 3
	Print               // 4
	Input               // 5
	StartLoop           // 6 // Special
	EndLoop             //   7

	// Operations that take an argument
	DataDecArg  // 8
	DataIncArg  // 9
	IndexDecArg //10
	IndexIncArg //11

	// Advanced ops
	Zero
	Move
	Copy
	Plus
	Minus
	Mult
	Exp
	Divide
)
