package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

const HEADC = `
#include <stdio.h>
#include <stdint.h>
#define BASESIZE 30000

int main() {
	//setbuf(stdout, NULL);

	char mem[BASESIZE] = { 0 };

	char *ptr = mem;
	int counter = 0;
`
const FOOTC = `
}
`
const optimization = "-O1"
const binaryName = "compiled"

func compileCAndExec(op Executable, startTime time.Time) (compileTime time.Duration) {
	const brainfuckCFile = "brainfuck.c"

	// Translate program to C code
	c := programToC(op)
	ioutil.WriteFile(brainfuckCFile, c, 0733)

	// Compile C code
	gcc := exec.Command("gcc", brainfuckCFile, optimization, "-o", binaryName)
	gcc.Stdout = os.Stdout
	gcc.Stderr = os.Stderr
	err := gcc.Run()
	if err != nil {
		println(err.Error())
	}

	os.Rename(brainfuckCFile, brainfuckCFile+".txt") // To avoid borking the Go compiler
	compileTime = time.Since(startTime)

	// Run compiled C code
	program := exec.Command("./" + binaryName)
	program.Stdout = os.Stdout
	program.Stderr = os.Stderr
	program.Run()
	if err != nil {
		println(err.Error())
	}
	return
}

func programToC(e Executable) []byte {
	var b bytes.Buffer

	// Write head of the file
	b.WriteString(HEADC)

	// Compile all instructions
	e.toC(&b)

	// Write last part of the file
	b.WriteString(FOOTC)

	return b.Bytes()
}
