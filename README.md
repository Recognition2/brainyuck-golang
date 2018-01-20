# Brainyuck parser, written in Go
[![Build Status](https://travis-ci.org/Recognition2/brainyuck-golang.svg?branch=master)](https://travis-ci.org/Recognition2/brainyuck-golang)

### How to work with it
The filename of the brainyuck code can be passed with the `-filename <file>` flag. Note the single dash (go default..)

This parser supports both interpret mode, enabled by passing the `-compile no` flag (default).

Alternatively, it is possible to compile the parsed bf code to C, and execute this C. On the mandelbrot benchmark, this gives ~16x higher performance.
This can be enabled with `-compile c`.


### Inner workings
The brainyuck code is first translated to intermediate instructions.
If interpreted mode is enabled, then these intermediate instructions are executed directly.
Otherwise, they are translated to a language of your choice (default C because it's fastest), and then compiled and executed.

### File structuring
All intermediate instructions implement the `Executable` interface, found in `executable.go`.
Possible instructions are `op`, `opWithArg`, `opWithArgOffset`, and `loop`. Their implementation can be found in the respective files.
A complete list of all supported instructions is given in `op.go`.

In interpret mode, all operations execute on a `State` struct. This struct contains the large array and the pointer that are typical for any brainyuck program. Furthermore it contains the output of the program (if buffering is enabled), and some statistics.

In compile mode:
- Advanced instructions are translated to their C equivalents and stored in a string,
- This string is written to a file,
- `gcc` is called on this file,
- The resulting binary is executed.


### How to use
```
go get github.com/recognition2/brainyuck-golang
$GOBIN/brainyuck-golang mandelbrot.bf
```

If you find a bug, please tell me.

Features:
- [x] Collapse sequential instructions (`>>>` becomes `>3`);
- [x] Collapse Zero loops (`[-]` or `[+]` nulls the current cell);
- [x] Collapse Seek loops (`[>>>]` means "increment by x until a cell with value zero is found)
- [x] Collapse loops consisting of only `><+-` instructions. If the pointer ends at the same value that it started with:
- [x] Collapse multiplication loops (called AddAndZero);

- [ ] Collapse nested loops (at the moment, only the most nested loops are collapsed where possible)
- [ ] Improve execution time of Go code (perhaps by reducing allocations?)
- [ ] Improve execution time EVEN MORE


Execution languages
- [x] Interpreted (Go)
- [x] Compiled as C
- [ ] Compiled as Go
- [ ] Compiled as Brainyuck (Bootstrapping)


### Why C:
Because it's trivial to call `gcc` on any file, it doesn't care about formatting, unused imports or other stuff, and does not know the concept of "borrowing".
This makes it the ideal language to generate.



