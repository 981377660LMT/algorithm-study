package main

import "fmt"

func main() {
	vm.Init()

	vm.Interpret(chunk, true)
}

// #region Chunk

type OpCode = byte // type of bytecode instructions

const (
	OP_CONSTANT OpCode = iota // load the constant for use
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
	OP_RETURN // return from the current function
)

type Chunk struct {
	code      []byte    // opcodes or operands
	constants []float64 // a chunk may only contain up to 256 different constants

	lines []int
}

func NewChunk() *Chunk {
	return &Chunk{}
}

func (c *Chunk) Write(v byte, line int) {
	c.code = append(c.code, v)
	c.lines = append(c.lines, line)
}

func (c *Chunk) AddConstant(v float64) byte {
	c.constants = append(c.constants, v)
	return byte(len(c.constants) - 1)
}

func (c *Chunk) Disassemble(name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < len(c.code); {
		offset = c.disassembleInstruction(offset)
	}
}

// Instructions can have different sizes.
func (c *Chunk) disassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)

	// show a | for any instruction that comes from the same source line as the preceding one
	if offset > 0 && c.lines[offset] == c.lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", c.lines[offset])
	}

	instruction := c.code[offset]
	switch instruction {
	case OP_CONSTANT:
		return c.constantInstruction("OP_CONSTANT", offset)
	case OP_ADD:
		return c.simpleInstruction("OP_ADD", offset)
	case OP_SUBTRACT:
		return c.simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		return c.simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return c.simpleInstruction("OP_DIVIDE", offset)
	case OP_NEGATE:
		return c.simpleInstruction("OP_NEGATE", offset)
	case OP_RETURN:
		return c.simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func (c *Chunk) constantInstruction(name string, offset int) int {
	constant := c.code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	fmt.Printf("%g", c.constants[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func (c *Chunk) simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

// #endregion

// #region VM

type InterpretResult byte

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

const STACK_MAX int = 256

var vm = NewVM() // use global variable to keep the code in the book a little lighter

type VM struct {
	chunk *Chunk

	// instruction pointer, keeps track of the current instruction being executed
	// the IP always points to `the next instruction`, not the one currently being handled
	ip int

	stack    [STACK_MAX]float64
	stackTop int // points to where the next value to be pushed will go
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) Init() {
}

func (vm *VM) Interpret(source string, debug bool) InterpretResult {
	// TODO: inject
	cp := NewCompiler(source)
	cp.compile()
	return vm.run(debug)
}

func (vm *VM) run(debug bool) InterpretResult {
	if debug {
		fmt.Println("         ")
		for i := 0; i < vm.stackTop; i++ {
			fmt.Printf("[ %g ]\n", vm.stack[i])
		}
		vm.chunk.Disassemble("test chunk")
	}

	for {
		// decoding/dispatching
		instruction := vm.readByte()
		switch instruction {
		case OP_CONSTANT:
			constant := vm.readConstant()
			vm.push(constant)
			break
		case OP_ADD:
			vm.binaryOp(func(a, b float64) float64 { return a + b })
			break
		case OP_SUBTRACT:
			vm.binaryOp(func(a, b float64) float64 { return a - b })
			break
		case OP_MULTIPLY:
			vm.binaryOp(func(a, b float64) float64 { return a * b })
			break
		case OP_DIVIDE:
			vm.binaryOp(func(a, b float64) float64 { return a / b })
			break
		case OP_NEGATE:
			vm.push(-vm.pop())
			break
		case OP_RETURN:
			fmt.Printf("%g\n", vm.pop())
			return INTERPRET_OK
		}
	}
}

func (vm *VM) readByte() byte {
	vm.ip++
	return vm.chunk.code[vm.ip-1]
}

func (vm *VM) readConstant() float64 {
	return vm.chunk.constants[vm.readByte()]
}

func (vm *VM) push(v float64) {
	vm.stack[vm.stackTop] = v
	vm.stackTop++
}

func (vm *VM) pop() float64 {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) binaryOp(f func(float64, float64) float64) {
	b := vm.pop()
	a := vm.pop()
	vm.push(f(a, b))
}

// #endregion

// #region Compiler
type Compiler struct {
	source string
}

func NewCompiler(source string) *Compiler {
	return &Compiler{source: source}
}

func (c *Compiler) compile() *Chunk {
	// TODO: inject
	scanner := NewScanner(c.source)
}

// #endregion

// #region Scanner
type Scanner struct {
	start   int // start of lexeme
	current int // current character of lexeme
	line    int // current line of lexeme

	source string
}

func NewScanner(source string) *Scanner {
	return &Scanner{line: 1, source: source}
}

// #endregion
