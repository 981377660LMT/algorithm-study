// !Source Code -> Scanner -> Tokens -> Parser、Compiler -> Bytecode Chunk -> VM -> Output

package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	// TestScanner("print 1 + 2;")

	// Run("3 + 4 / 5") // expect 23
	// Run("true")
	// !(5 - 4 > 3 * 2 == !nil)
	// Run("print !(5 - 4 > 3 * 2 == !nil);")
	// Run(`print 1+2;`)

	// var beverage = "cafe au lait";
	// var breakfast = "beignets with " + beverage;
	// print breakfast;

	var sources string

	sources = `
	1 + 2;
	`
	Run(sources)

	sources = `
	{
		var a = "outer";
		{
		 var a = "inner";
		 print a;
		}
		print a;
	}`
	Run(sources)

	// sources = `
	// var a = 1;
	// if (a > 0) {
	// 	print "positive";
	// } else {
	// 	print "negative";
	// }`
	// Run(sources)

	// sources = `
	// print (1 and 2);
	// print (1 or 2);
	// `
	// Run(sources)

	// sources = `
	// var a = 1;
	// while (a < 10) {
	// 	print a;
	// 	a = a + 1;
	// }
	// `
	// Run(sources)

	// // for
	// sources = `
	// for (var i = 0; i < 10; i = i + 1) {
	// 	print i;
	// }
	// `
	// Run(sources)
}

func TestScanner(source string) {
	S := NewScanner(source)
	S.Scan()
}

func Run(source string) {
	scanner := NewScanner(source)
	resolver := NewResolver(TYPE_SCRIPT)
	compiler := NewCompiler(scanner, resolver)

	vm.Init(compiler)
	vm.Interpret(source, true)
	vm.Free()
}

// #region Chunk 字节码块，供VM执行

type OpCode = byte // type of bytecode instructions

const (
	OP_CONSTANT OpCode = iota // load the constant for use
	OP_NIL
	OP_TRUE
	OP_FALSE
	OP_POP
	OP_GET_LOCAL
	OP_SET_LOCAL
	OP_GET_GLOBAL
	OP_DEFINE_GLOBAL
	OP_SET_GLOBAL
	OP_EQUAL
	OP_GREATER
	OP_LESS
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NOT
	OP_NEGATE
	OP_PRINT
	OP_JUMP
	OP_JUMP_IF_FALSE
	OP_LOOP
	OP_RETURN // return from the current function
)

// 字节码是一系列指令.充当jlox中AST的作用.
// 可以将字节码视为 AST 的一种紧凑序列化.
type Chunk struct {
	code      []byte   // opcodes or operands
	constants []IValue // 常量表，a chunk may only contain up to 256 different constants

	lines []int
}

func NewChunk() *Chunk {
	return &Chunk{}
}

func (c *Chunk) String() string {
	return fmt.Sprintf("Chunk{code: %v, constants: %v, lines: %v}", c.code, c.constants, c.lines)
}

func (c *Chunk) Write(b byte, line int) {
	c.code = append(c.code, b)
	c.lines = append(c.lines, line)
}

func (c *Chunk) AddConstant(v IValue) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
}

// 反汇编器将CPU指令转换为人类可读的指令
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
	case OP_NIL:
		return c.simpleInstruction("OP_NIL", offset)
	case OP_TRUE:
		return c.simpleInstruction("OP_TRUE", offset)
	case OP_FALSE:
		return c.simpleInstruction("OP_FALSE", offset)
	case OP_POP:
		return c.simpleInstruction("OP_POP", offset)
	case OP_GET_LOCAL:
		return c.byteInstruction("OP_GET_LOCAL", offset)
	case OP_SET_LOCAL:
		return c.byteInstruction("OP_SET_LOCAL", offset)
	case OP_GET_GLOBAL:
		return c.constantInstruction("OP_GET_GLOBAL", offset)
	case OP_DEFINE_GLOBAL:
		return c.constantInstruction("OP_DEFINE_GLOBAL", offset)
	case OP_SET_GLOBAL:
		return c.constantInstruction("OP_SET_GLOBAL", offset)
	case OP_EQUAL:
		return c.simpleInstruction("OP_EQUAL", offset)
	case OP_GREATER:
		return c.simpleInstruction("OP_GREATER", offset)
	case OP_LESS:
		return c.simpleInstruction("OP_LESS", offset)
	case OP_ADD:
		return c.simpleInstruction("OP_ADD", offset)
	case OP_SUBTRACT:
		return c.simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		return c.simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return c.simpleInstruction("OP_DIVIDE", offset)
	case OP_NOT:
		return c.simpleInstruction("OP_NOT", offset)
	case OP_NEGATE:
		return c.simpleInstruction("OP_NEGATE", offset)
	case OP_PRINT:
		return c.simpleInstruction("OP_PRINT", offset)
	case OP_JUMP:
		return c.jumpInstruction("OP_JUMP", 1, offset)
	case OP_JUMP_IF_FALSE:
		return c.jumpInstruction("OP_JUMP_IF_FALSE", 1, offset)
	case OP_LOOP:
		return c.jumpInstruction("OP_LOOP", -1, offset)
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
	fmt.Printf("%v", c.constants[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func (c *Chunk) simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func (c *Chunk) byteInstruction(name string, offset int) int {
	slot := c.code[offset+1]
	fmt.Printf("%-16s %4d\n", name, slot)
	return offset + 2
}

func (c *Chunk) jumpInstruction(name string, sign int, offset int) int {
	jump := int(c.code[offset+1])<<8 | int(c.code[offset+2])
	fmt.Printf("%-16s %4d -> %d\n", name, offset, offset+3+sign*jump)
	return offset + 3
}

// #endregion

// #region Obj 对象，存储在堆上
type ObjType byte

const (
	OBJ_STRING ObjType = iota
	OBJ_FUNCTION
)

type Obj interface {
	Type() ObjType
	Value() any

	// 最简单的方法跟踪对象，用于垃圾回收.
	Next() Obj
	HashCode() int
	String() string
}

type FunctionType byte // 让编译器能够区分它是在编译顶层代码还是函数体

const (
	TYPE_FUNCTION FunctionType = iota
	TYPE_SCRIPT
)

// Factory function for creating objects.
func NewObj(t ObjType, args ...any) Obj {
	switch t {
	case OBJ_STRING:
		return NewObjString(args[0].(string))
	case OBJ_FUNCTION:
		return NewObjFunction(args[0].(string), args[1].(int), args[2].(*Chunk))
	}
	return nil
}

type ObjString struct {
	next Obj

	v string
}

func NewObjString(v string) *ObjString {
	res := &ObjString{v: v}
	res.next = vm.objects
	vm.objects = res
	return res
}
func (o *ObjString) Next() Obj      { return o.next }
func (o *ObjString) Type() ObjType  { return OBJ_STRING }
func (o *ObjString) Value() any     { return o.v }
func (o *ObjString) HashCode() int  { return hashString(o.v) }
func (o *ObjString) String() string { return o.v }

type ObjFunction struct {
	next Obj

	name  string
	arity int
	chunk *Chunk
}

func NewObjFunction(name string, arity int, chunk *Chunk) *ObjFunction {
	res := &ObjFunction{name: name, arity: arity, chunk: chunk}
	res.next = vm.objects
	vm.objects = res
	return res
}

func (o *ObjFunction) Next() Obj     { return o.next }
func (o *ObjFunction) Type() ObjType { return OBJ_FUNCTION }
func (o *ObjFunction) Value() any    { return o.chunk }
func (o *ObjFunction) HashCode() int { return hashString(o.name) }
func (o *ObjFunction) String() string {
	if o.name == "" {
		return "<script>"
	}
	return fmt.Sprintf("<fn %s>", o.name)
}

// #endregion

// #region Value 类型系统相关
type ValueType byte

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER

	VAL_OBJ
)

// 包装类型.
type IValue interface {
	Type() ValueType
	Value() any

	ToBool() bool
	ToNumber() float64

	HashCode() int
}

type BoolValue struct {
	v bool
}

func NewBoolValue(v bool) IValue     { return &BoolValue{v} }
func (v *BoolValue) Type() ValueType { return VAL_BOOL }
func (v *BoolValue) Value() any      { return v.v }
func (v *BoolValue) ToBool() bool    { return v.v }
func (v *BoolValue) ToNumber() float64 {
	if v.v {
		return 1
	}
	return 0
}
func (v *BoolValue) HashCode() int {
	if v.v {
		return 1231
	}
	return 1237
}
func (v *BoolValue) String() string {
	if v.v {
		return "true"
	}
	return "false"
}

type NilValue struct{}

func NewNilValue() IValue {
	return &NilValue{}
}
func (v *NilValue) Type() ValueType   { return VAL_NIL }
func (v *NilValue) Value() any        { return nil }
func (v *NilValue) ToBool() bool      { return false }
func (v *NilValue) ToNumber() float64 { return 0 }
func (v *NilValue) HashCode() int     { return 0 }
func (v *NilValue) String() string    { return "nil" }

type NumberValue struct {
	v float64
}

func NewNumberValue(v float64) IValue    { return &NumberValue{v} }
func (v *NumberValue) Type() ValueType   { return VAL_NUMBER }
func (v *NumberValue) Value() any        { return v.v }
func (v *NumberValue) ToBool() bool      { return v.v != 0 }
func (v *NumberValue) ToNumber() float64 { return v.v }
func (v *NumberValue) HashCode() int     { return int(v.v) }
func (v *NumberValue) String() string    { return fmt.Sprintf("%g", v.v) }

type ObjValue struct {
	v Obj
}

func NewValueObj(v Obj) IValue { return &ObjValue{v} }
func (v *ObjValue) Type() ValueType {
	return VAL_OBJ
}
func (v *ObjValue) Value() any        { return v.v }
func (v *ObjValue) ToBool() bool      { return true }
func (v *ObjValue) ToNumber() float64 { return 0 }
func (v *ObjValue) HashCode() int     { return v.v.HashCode() }
func (v *ObjValue) String() string    { return v.v.String() }

func IsSameValue(a, b IValue) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch a.Type() {
	case VAL_BOOL:
		return a.ToBool() == b.ToBool()
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return a.ToNumber() == b.ToNumber()
	case VAL_OBJ:
		return a.HashCode() == b.HashCode()
	default:
		return false
	}
}
func IsBool(v IValue) bool   { return v.Type() == VAL_BOOL }
func IsNil(v IValue) bool    { return v.Type() == VAL_NIL }
func IsNumber(v IValue) bool { return v.Type() == VAL_NUMBER }
func IsObj(v IValue) bool    { return v.Type() == VAL_OBJ }
func IsStringObj(v IValue) bool {
	if IsObj(v) {
		return v.Value().(Obj).Type() == OBJ_STRING
	}
	return false
}
func IsFunctionObj(v IValue) bool {
	if IsObj(v) {
		return v.Value().(Obj).Type() == OBJ_FUNCTION
	}
	return false
}

// #endregion

// #region VM 虚拟机

type InterpretResult byte

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

const FRAMES_MAX int = 64
const UINT8_COUNT int = 256
const STACK_MAX int = FRAMES_MAX * UINT8_COUNT

// 如果要向所有函数传递一个指向虚拟机的指针，会很麻烦
// use global variable to keep the code in the book a little lighter
var vm = NewVM()

// 每次调用函数时，我们都会创建一个新的 CallFrame 并将其推送到虚拟机的帧堆栈中。
type CallFrame struct {
	function *ObjFunction // 被调用函数
	ip       int          // 返回地址。从一个函数返回时，虚拟机将跳转到调用者的 CallFrame 的 ip 并从那里恢复执行。
	base     int          // 虚拟机值栈中，该函数可以使用的第一个槽的索引
}

func NewCallFrame(function *ObjFunction, ip int, slotsStart int) *CallFrame {
	return &CallFrame{function: function, ip: ip, base: slotsStart}
}

func (frame *CallFrame) String() string {
	return fmt.Sprintf("CallFrame{function: %s, ip: %d, slotsStart: %d}", frame.function, frame.ip, frame.base)
}

type VM struct {
	// chunk *Chunk

	// // 指令指针/程序计数器，用于跟踪当前正在执行的指令。我们在程序中的“位置”。
	// // instruction pointer, keeps track of the current instruction being executed
	// // the IP always points to `the next instruction`, not the one currently being handled
	// // chunk.code[ip]
	// ip int

	frames     [FRAMES_MAX]*CallFrame // 对应vm.stack的一个窗口
	frameCount int

	stack    [STACK_MAX]IValue
	stackTop int // points to where the next value to be pushed will go

	objects Obj
	globals map[int]IValue // 全局变量表

	// TODO: strings pool

	compiler *Compiler
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) Init(compiler *Compiler) {
	vm.compiler = compiler
	for i := 0; i < FRAMES_MAX; i++ {
		vm.frames[i] = NewCallFrame(nil, 0, 0)
	}
	vm.frameCount = 0
	vm.stackTop = 0
	vm.objects = nil
	vm.globals = make(map[int]IValue)
}

func (vm *VM) Free() {
	ptr := vm.objects
	for ptr != nil {
		next := ptr.Next()
		vm.free(ptr)
		ptr = next
	}
}

// 想象这里有一个垃圾回收器，它会在这里运行。
func (vm *VM) free(obj Obj) {
}

func (vm *VM) Interpret(source string, debug bool) InterpretResult {
	function := vm.compiler.Compile(source)
	if function == nil {
		return INTERPRET_COMPILE_ERROR
	}

	vm.push(NewValueObj(function))
	frame := vm.frames[vm.frameCount]
	vm.frameCount++
	frame.function = function
	frame.ip = 0
	frame.base = 0
	fmt.Println(frame.function.chunk)

	res := vm.run(debug)
	return res
}

func (vm *VM) run(debug bool) InterpretResult {
	allowDisassemble := debug && !vm.compiler.hadError
	if allowDisassemble {
		fmt.Println("         ")
		for i := 0; i < vm.stackTop; i++ {
			fmt.Printf("[ %s ]\n", vm.stack[i])
		}
		frame := vm.frame()
		var name string
		if frame.function.name == "" {
			name = "<script>"
		} else {
			name = frame.function.name
		}
		frame.function.chunk.Disassemble(name)
	}

	for {
		// decoding/dispatching
		instruction := vm.readByte()
		switch instruction {
		case OP_CONSTANT:
			vm.push(vm.readConstant())
			break
		case OP_NIL:
			vm.push(NewNilValue())
			break
		case OP_TRUE:
			vm.push(NewBoolValue(true))
			break
		case OP_FALSE:
			vm.push(NewBoolValue(false))
			break
		case OP_POP:
			vm.pop()
			break
		case OP_GET_LOCAL:
			slot := int(vm.readByte())
			// 访问当前帧的 slots 数组
			vm.push(*vm.slotAt(slot))
			break
		case OP_SET_LOCAL:
			slot := int(vm.readByte())
			*vm.slotAt(slot) = vm.peek(0)
			// Don't pop, since the set operation has the RHS as its return value.
			break
		case OP_GET_GLOBAL:
			name := vm.readConstant()
			if v, ok := vm.globals[name.HashCode()]; !ok {
				vm.runtimeError(fmt.Sprintf("Undefined variable '%s'.", name))
				return INTERPRET_RUNTIME_ERROR
			} else {
				vm.push(v)
			}
			break
		case OP_DEFINE_GLOBAL:
			// 从常量表中获取变量的名称，然后我们从栈顶获取值，并以该名称为键将其存储在哈希表中
			name := vm.readConstant()
			vm.globals[name.HashCode()] = vm.pop()
			break
		case OP_SET_GLOBAL:
			name := vm.readConstant()
			hashCode := name.HashCode()
			if _, ok := vm.globals[hashCode]; !ok {
				vm.runtimeError(fmt.Sprintf("Undefined variable '%s'.", name))
				return INTERPRET_RUNTIME_ERROR
			}
			vm.globals[hashCode] = vm.peek(0)
			break
		case OP_EQUAL:
			b := vm.pop()
			a := vm.pop()
			vm.push(NewBoolValue(IsSameValue(a, b)))
			break
		case OP_GREATER:
			vm.binaryOp(func(a, b float64) IValue { return NewBoolValue(a > b) })
			break
		case OP_LESS:
			vm.binaryOp(func(a, b float64) IValue { return NewBoolValue(a < b) })
			break
		case OP_ADD:
			if IsNumber(vm.peek(0)) && IsNumber(vm.peek(1)) {
				vm.binaryOp(func(a, b float64) IValue { return NewNumberValue(a + b) })
			} else if IsStringObj(vm.peek(0)) && IsStringObj(vm.peek(1)) {
				vm.concatenate()
			} else {
				vm.runtimeError("Operands must be two numbers or two strings.")
				return INTERPRET_RUNTIME_ERROR
			}
			break
		case OP_SUBTRACT:
			vm.binaryOp(func(a, b float64) IValue { return NewNumberValue(a - b) })
			break
		case OP_MULTIPLY:
			vm.binaryOp(func(a, b float64) IValue { return NewNumberValue(a * b) })
			break
		case OP_DIVIDE:
			vm.binaryOp(func(a, b float64) IValue { return NewNumberValue(a / b) })
			break
		case OP_NOT:
			vm.push(NewBoolValue(vm.isFalsey(vm.pop())))
			break
		case OP_NEGATE:
			if !IsNumber(vm.peek(0)) {
				vm.runtimeError("Operand must be a number.")
				return INTERPRET_RUNTIME_ERROR
			}
			vm.push(NewNumberValue(-vm.pop().ToNumber()))
			break
		case OP_PRINT:
			fmt.Println(vm.pop())
			break
		case OP_JUMP:
			jumpLen := vm.readShort()
			vm.frame().ip += jumpLen
			break
		case OP_JUMP_IF_FALSE:
			jumpLen := vm.readShort()
			if vm.isFalsey(vm.peek(0)) {
				vm.frame().ip += jumpLen
			}
			break
		case OP_LOOP:
			jumpLen := vm.readShort()
			vm.frame().ip -= jumpLen // loop back
			break
		case OP_RETURN:
			return INTERPRET_OK
		}
	}
}

func (vm *VM) readByte() (res byte) {
	res = vm.chunk().code[*vm.ip()]
	*vm.ip()++
	return
}

func (vm *VM) readShort() (res int) {
	res = int(vm.readByte()) << 8
	res |= int(vm.readByte())
	return
}

func (vm *VM) readConstant() (res IValue) {
	res = vm.chunk().constants[vm.readByte()]
	return
}

func (vm *VM) push(v IValue) {
	vm.stack[vm.stackTop] = v
	vm.stackTop++
}

func (vm *VM) pop() IValue {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

// 将操作数留在栈上是很重要的，可以确保在运行过程中触发垃圾收集时，垃圾收集器能够找到它们.
func (vm *VM) peek(distance int) IValue {
	return vm.stack[vm.stackTop-1-distance]
}

func (vm *VM) isFalsey(v IValue) bool {
	return IsNil(v) || (IsBool(v) && !v.ToBool())
}

func (vm *VM) concatenate() {
	b := vm.pop().Value().(*ObjString).v
	a := vm.pop().Value().(*ObjString).v
	vm.push(NewValueObj(NewObjString(a[:len(a)-1] + b[1:])))
}

func (vm *VM) binaryOp(f func(float64, float64) IValue) {
	if !IsNumber(vm.peek(0)) || !IsNumber(vm.peek(1)) {
		vm.runtimeError("Operands must be numbers.")
		return
	}
	b := (vm.pop().ToNumber())
	a := (vm.pop().ToNumber())
	vm.push(f(a, b))
}

func (vm *VM) runtimeError(message string) {
	fmt.Printf("%s\n", message)

	// 从栈顶的调用帧中提取信息
	line := vm.chunk().lines[*vm.ip()]
	fmt.Printf("[line %d] in script\n", line)
}

func (vm *VM) frame() *CallFrame {
	if vm.frameCount == 0 {
		return nil
	}
	return vm.frames[vm.frameCount-1]
}

func (vm *VM) ip() *int {
	frame := vm.frame()
	if frame == nil {
		return nil
	}
	return &frame.ip
}

func (vm *VM) chunk() *Chunk {
	frame := vm.frame()
	if frame == nil || frame.function == nil {
		return nil
	}
	return frame.function.chunk
}

func (vm *VM) slotAt(idx int) *IValue {
	idx1 := vm.slotIdxAt(idx)
	if idx1 < 0 || idx1 >= vm.stackTop {
		return nil
	}
	return &vm.stack[idx1]
}

func (vm *VM) slotIdxAt(idx int) (stackIdx int) {
	frame := vm.frame()
	if frame == nil {
		return -1
	}
	return frame.base + idx
}

// #endregion

// #region Compiler 编译器产生字节码块

// Pratt parser 算法.
// 函数只解析一种类型的表达式。它们不会级联以包含更高优先级的表达式类型。
type Compiler struct {
	*state
	rules []*ParseRule

	scanner  *Scanner
	resolver *Resolver
}

func NewCompiler(scanner *Scanner, resolver *Resolver) *Compiler {
	res := &Compiler{state: NewState(), scanner: scanner, resolver: resolver}
	res.rules = res.createRules()
	return res
}

// 解析源代码并输出低级的二进制指令序列。
// 当然，它是字节码，而不是某个芯片的原生指令集，但它比jlox更接近于硬件。
// 接收源代码，返回一个包含字节码的函数对象。
func (c *Compiler) Compile(source string) *ObjFunction {
	c.advance()
	for !c.match(TOKEN_EOF) {
		c.declaration()
	}
	function := c.endCompiler()
	if c.hadError {
		return nil
	}
	return function
}

func (c *Compiler) advance() {
	c.previous = c.current
	for {
		c.current = c.scanner.ScanToken()
		if c.current.typ != TOKEN_ERROR {
			break
		}
		c.errorAtCurrent(c.current.value)
	}
}

func (c *Compiler) expression() {
	// 我们只需要解析最低优先级，它也包含了所有更高优先级的表达式。
	c.parsePrecedence(PREC_ASSIGNMENT)
}

func (c *Compiler) varDeclaration() {
	global := c.parseVariable("Expect variable name.")
	if c.match(TOKEN_EQUAL) {
		c.expression()
	} else {
		c.emitByte(OP_NIL)
	}
	c.consume(TOKEN_SEMICOLON, "Expect ';' after variable declaration.")
	c.defineVariable(global)
}

// block -> "{" declaration* "}"
func (c *Compiler) block() {
	for !c.check(TOKEN_RIGHT_BRACE) && !c.check(TOKEN_EOF) {
		c.declaration()
	}
	c.consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.")
}

// “表达式语句”就是一个表达式后面跟着一个分号。
// 从语义上说，表达式语句会对表达式求值并丢弃结果。
func (c *Compiler) expressionStatement() {
	c.expression()
	c.consume(TOKEN_SEMICOLON, "Expect ';' after expression.")
	c.emitByte(OP_POP)
}

// forStatement ->
//
//	"for" "(" ( var | expression? ";" expression? ";" expression? ) ")" statement
func (c *Compiler) forStatement() {
	c.beginScope() // for 块级作用域

	c.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'for'.")
	if c.match(TOKEN_SEMICOLON) {
		// No initializer.
	} else if c.match(TOKEN_VAR) {
		c.varDeclaration()
	} else {
		c.expressionStatement()
	}

	loopStart := len(c.currentChunk().code)
	exitJump := -1 // 退出循环的条件表达式
	if !c.match(TOKEN_SEMICOLON) {
		c.expression()
		c.consume(TOKEN_SEMICOLON, "Expect ';' after loop condition.")

		// Jump out of the loop if the condition is false.
		exitJump = c.emitJump(OP_JUMP_IF_FALSE)
		c.emitByte(OP_POP) // Condition.
	}

	if !c.match(TOKEN_RIGHT_PAREN) {
		bodyJump := c.emitJump(OP_JUMP)
		incrementStart := len(c.currentChunk().code) // 记录增量表达式的起始位置
		c.expression()
		c.emitByte(OP_POP)
		c.consume(TOKEN_RIGHT_PAREN, "Expect ')' after for clauses.")

		c.emitLoop(loopStart)
		loopStart = incrementStart
		c.patchJump(bodyJump)
	}

	c.statement()
	c.emitLoop(loopStart)

	if exitJump != -1 {
		c.patchJump(exitJump)
		c.emitByte(OP_POP) // Condition.
	}

	c.endScope()
}

// ifStatement -> "if" "(" expression ")" statement ("else" statement)?
func (c *Compiler) ifStatement() {
	c.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'if'.")
	c.expression()
	c.consume(TOKEN_RIGHT_PAREN, "Expect ')' after condition.")

	thenJump := c.emitJump(OP_JUMP_IF_FALSE)
	c.emitByte(OP_POP)
	c.statement()
	elseJump := c.emitJump(OP_JUMP)
	c.patchJump(thenJump)
	c.emitByte(OP_POP)
	if c.match(TOKEN_ELSE) {
		c.statement()
	}
	c.patchJump(elseJump)
}

func (c *Compiler) printStatement() {
	c.expression()
	c.consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	c.emitByte(OP_PRINT)
}

// whileStatement -> "while" "(" expression ")" statement
func (c *Compiler) whileStatement() {
	loopStart := len(c.currentChunk().code)
	c.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'while'.")
	c.expression()
	c.consume(TOKEN_RIGHT_PAREN, "Expect ')' after condition.")

	exitJump := c.emitJump(OP_JUMP_IF_FALSE)
	c.emitByte(OP_POP)
	c.statement()
	c.emitLoop(loopStart)

	c.patchJump(exitJump)
	c.emitByte(OP_POP)
}

func (c *Compiler) synchronize() {
	c.panicMode = false
	for c.current.typ != TOKEN_EOF {
		if c.previous.typ == TOKEN_SEMICOLON {
			return
		}
		switch c.current.typ {
		case TOKEN_CLASS, TOKEN_FUN, TOKEN_VAR, TOKEN_FOR, TOKEN_IF, TOKEN_WHILE, TOKEN_PRINT, TOKEN_RETURN:
			return
		}
		c.advance()
	}
}

func (c *Compiler) declaration() {
	if c.match(TOKEN_VAR) {
		c.varDeclaration()
	} else {
		c.statement()
	}

	if c.panicMode {
		c.synchronize()
	}
}

func (c *Compiler) statement() {
	if c.match(TOKEN_PRINT) {
		c.printStatement()
	} else if c.match(TOKEN_FOR) {
		c.forStatement()
	} else if c.match(TOKEN_IF) {
		c.ifStatement()
	} else if c.match(TOKEN_WHILE) {
		c.whileStatement()
	} else if c.match(TOKEN_LEFT_BRACE) {
		c.beginScope()
		c.block()
		c.endScope()
	} else {
		c.expressionStatement()
	}
}

// 读取下一个标识、验证标识是否具有预期的类型。如果不是，则报告错误。
func (c *Compiler) consume(tokenType TokenType, message string) {
	if c.current.typ == tokenType {
		c.advance()
		return
	}
	c.errorAtCurrent(message)
}

func (c *Compiler) check(tokenType TokenType) bool {
	return c.current.typ == tokenType
}

func (c *Compiler) match(tokenType TokenType) bool {
	if !c.check(tokenType) {
		return false
	}
	c.advance()
	return true
}

// 向块中追加一个字节。
// 将给定的字节写入一个指令，该字节可以是操作码或操作数。
// 它会发送前一个token的行信息，以便将运行时错误与该行关联起来。
func (c *Compiler) emitByte(b byte) {
	c.currentChunk().Write(b, c.previous.line)
}

func (c *Compiler) emitLoop(loopStart int) {
	c.emitByte(OP_LOOP)
	// 从当前指令计算到我们想要跳回的 loopStart 点的偏移量。
	// + 2 是考虑到 OP_LOOP 指令自身操作数的大小，我们也需要跳过这些操作数。
	offset := len(c.currentChunk().code) - loopStart + 2
	if offset > math.MaxUint16 {
		c.error("Loop body too large.")
	}
	c.emitByte(byte((offset >> 8) & 0xff))
	c.emitByte(byte(offset & 0xff))
}

// 一般是写一个操作码，后面跟一个单字节的操作数。
func (c *Compiler) emitByte2(b1, b2 byte) {
	c.emitByte(b1)
	c.emitByte(b2)
}

func (c *Compiler) emitReturn() {
	c.emitByte(OP_RETURN)
}

func (c *Compiler) makeConstant(value IValue) byte {
	constant := c.currentChunk().AddConstant(value)
	// OP_CONSTANT指令使用单个字节来索引操作数，所以我们在一个块中最多只能存储和加载256个常量。
	if constant > math.MaxUint8 {
		c.error("Too many constants in one chunk")
		return 0
	}
	return byte(constant)
}

func (c *Compiler) emitConstant(value IValue) {
	c.emitByte2(OP_CONSTANT, c.makeConstant(value))
}

func (c *Compiler) emitJump(instruction OpCode) int {
	c.emitByte(instruction)
	// 16 位的偏移量让我们可以跳过最多 65,535 字节的代码
	c.emitByte(0xff)
	c.emitByte(0xff)
	return len(c.currentChunk().code) - 2
}

// 根据偏移量修补之前的emitJump指令.
func (c *Compiler) patchJump(offset int) {
	jumpLen := len(c.currentChunk().code) - offset - 2 // 字节码长度
	if jumpLen > math.MaxUint16 {
		c.error("Too much code to jump over.")
	}
	c.currentChunk().code[offset] = byte((jumpLen >> 8) & 0xff)
	c.currentChunk().code[offset+1] = byte(jumpLen & 0xff)
}

func (c *Compiler) endCompiler() *ObjFunction {
	c.emitReturn()
	return c.resolver.function
}

// 进入一个新的作用域.
func (c *Compiler) beginScope() {
	c.resolver.scopeDepth++
}

func (c *Compiler) endScope() {
	current := c.resolver
	current.scopeDepth--

	// remove tail
	for current.localCount > 0 && current.locals[current.localCount-1].depth > current.scopeDepth {
		c.emitByte(OP_POP) // 可以优化成OP_POPN, 一条指令弹出多个值
		current.localCount--
	}
}

func (c *Compiler) binary(canAssign bool) {
	typ := c.previous.typ
	// 为每个二元运算符定义一个单独的函数。每个函数都会调用 parsePrecedence() 并传入正确的优先级来解析其操作数。
	rule := c.getRule(typ)
	c.parsePrecedence(rule.precedence + 1)

	switch typ {
	case TOKEN_BANG_EQUAL:
		c.emitByte2(OP_EQUAL, OP_NOT)
		break
	case TOKEN_EQUAL_EQUAL:
		c.emitByte(OP_EQUAL)
		break
	case TOKEN_GREATER:
		c.emitByte(OP_GREATER)
		break
	case TOKEN_GREATER_EQUAL:
		c.emitByte2(OP_LESS, OP_NOT)
		break
	case TOKEN_LESS:
		c.emitByte(OP_LESS)
		break
	case TOKEN_LESS_EQUAL:
		c.emitByte2(OP_GREATER, OP_NOT)
		break
	case TOKEN_PLUS:
		c.emitByte(OP_ADD)
		break
	case TOKEN_MINUS:
		c.emitByte(OP_SUBTRACT)
		break
	case TOKEN_STAR:
		c.emitByte(OP_MULTIPLY)
		break
	case TOKEN_SLASH:
		c.emitByte(OP_DIVIDE)
		break
	default:
		return // Unreachable.
	}
}

func (c *Compiler) literal(canAssign bool) {
	// 当解析器遇到 false、nil 或 true 时，在前缀位置，它调用这个新的解析器函数.
	switch c.previous.typ {
	case TOKEN_FALSE:
		c.emitByte(OP_FALSE)
		break
	case TOKEN_NIL:
		c.emitByte(OP_NIL)
		break
	case TOKEN_TRUE:
		c.emitByte(OP_TRUE)
		break
	default:
		return
	}
}

func (c *Compiler) grouping(canAssign bool) {
	// 假定初始的(已经被消耗了。递归地调用expression()来编译括号之间的表达式，然后解析结尾的)。
	c.expression()
	c.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (c *Compiler) number(canAssign bool) {
	num, _ := strconv.ParseFloat(c.previous.value, 64)
	c.emitConstant(NewNumberValue(num))
}

func (c *Compiler) string(canAssign bool) {
	c.emitConstant(NewValueObj(NewObjString(c.previous.value)))
}

// 变量arg的set/get.
func (c *Compiler) namedVariable(name *Token, canAssign bool) {
	var setOp, getOp OpCode
	arg := c.resolveLocal(name)
	if arg != -1 {
		getOp = OP_GET_LOCAL
		setOp = OP_SET_LOCAL
	} else {
		arg = int(c.identifierConstant(name))
		getOp = OP_GET_GLOBAL
		setOp = OP_SET_GLOBAL
	}

	// 赋值, 例如: var a = 1;
	if canAssign && c.match(TOKEN_EQUAL) {
		c.expression()
		c.emitByte2(setOp, byte(arg))
	} else {
		// 取值, 例如: print a;
		c.emitByte2(getOp, byte(arg))
	}
}

func (c *Compiler) variable(canAssign bool) {
	c.namedVariable(c.previous, canAssign)
}

func (c *Compiler) unary(canAssign bool) {
	typ := c.previous.typ
	c.parsePrecedence(PREC_UNARY)
	switch typ {
	case TOKEN_BANG:
		c.emitByte(OP_NOT)
		break
	case TOKEN_MINUS:
		c.emitByte(OP_NEGATE)
		break
	default:
		return
	}
}

// Pratt parser 算法.
func (c *Compiler) parsePrecedence(p Precedence) {
	c.advance()
	prefixParselet := c.getRule(c.previous.typ).prefix
	if prefixParselet == nil {
		c.error("Expect expression.")
		return
	}
	canAssign := p <= PREC_ASSIGNMENT
	prefixParselet(canAssign)

	// 至此，整个左侧操作数表达式已经被编译，随后的中缀运算符也已被消耗.
	for p <= c.getRule(c.current.typ).precedence {
		c.advance()
		infixParselet := c.getRule(c.previous.typ).infix
		infixParselet(canAssign)
	}
	// 如果=没有作为表达式的一部分被消耗，那么其它任何东西都不会消耗它。
	if canAssign && c.match(TOKEN_EQUAL) {
		c.error("Invalid assignment target.")
	}
}

func (c *Compiler) defineVariable(global byte) {
	if c.resolver.scopeDepth > 0 {
		c.markInitialized()
		return
	}
	c.emitByte2(OP_DEFINE_GLOBAL, global)
}

// 在and表达式中，如果左侧为假值，则跳过右侧操作数.
func (c *Compiler) and_(canAssign bool) {
	endJump := c.emitJump(OP_JUMP_IF_FALSE)
	c.emitByte(OP_POP)
	c.parsePrecedence(PREC_AND)
	c.patchJump(endJump)
}

// 在or表达式中，如果左侧为真值，则跳过有操作数.
func (c *Compiler) or_(canAssign bool) {
	elseJump := c.emitJump(OP_JUMP_IF_FALSE)
	endJump := c.emitJump(OP_JUMP)
	c.patchJump(elseJump)
	c.emitByte(OP_POP)
	c.parsePrecedence(PREC_OR)
	c.patchJump(endJump)
}

func (c *Compiler) identifierConstant(name *Token) byte {
	return c.makeConstant(NewValueObj(NewObjString(name.value)))
}

func (c *Compiler) identifiersEqual(a, b *Token) bool {
	if a == nil || b == nil {
		return false
	}
	if a.typ != b.typ {
		return false
	}
	return a.value == b.value
}

// 返回局部变量在当前作用域中的索引.
// 如果变量尚未声明，则返回-1.
//
// !反向查找确保了内部本地变量正确地遮蔽了周围范围内同名的本地变量。
func (c *Compiler) resolveLocal(name *Token) int {
	for i := c.resolver.localCount - 1; i >= 0; i-- {
		local := c.resolver.locals[i]
		if c.identifiersEqual(name, local.name) {
			if local.depth == -1 {
				c.error("Cannot read local variable in its own initializer.")
			}
			return i
		}
	}
	return -1
}

// 关键函数：解析变量并将其添加到局部变量表中.
func (c *Compiler) addLocal(name *Token) {
	if c.resolver.localCount == 256 {
		c.error("Too many local variables in function.")
		return
	}
	local := c.resolver.locals[c.resolver.localCount]
	c.resolver.localCount++
	local.name = name
	local.depth = -1 // 未初始化(uninitialized)
}

func (c *Compiler) declareVariable() {
	if c.resolver.scopeDepth == 0 {
		return
	}
	name := c.previous

	// 检测此作用域中是否已经声明了同名变量(我们不允许Name Shadowing).
	for i := c.resolver.localCount - 1; i >= 0; i-- {
		local := c.resolver.locals[i]
		if local.depth != -1 && local.depth < c.resolver.scopeDepth { // pruning
			break
		}
		if c.identifiersEqual(name, local.name) {
			c.error("Variable with this name already declared in this scope.")
		}
	}

	c.addLocal(name)
}

// 返回该常量在常量表中的索引.
func (c *Compiler) parseVariable(errorMessage string) byte {
	c.consume(TOKEN_IDENTIFIER, errorMessage)

	c.declareVariable()
	if c.resolver.scopeDepth > 0 {
		return 0
	}

	return c.identifierConstant(c.previous)
}

func (c *Compiler) markInitialized() {
	c.resolver.locals[c.resolver.localCount-1].depth = c.resolver.scopeDepth
}

func (c *Compiler) currentChunk() *Chunk {
	return c.resolver.function.chunk
}

func (c *Compiler) errorAtCurrent(message string) {
	c.errorAt(c.current, message)
}

func (c *Compiler) error(message string) {
	c.errorAt(c.previous, message)
}

func (c *Compiler) errorAt(token *Token, message string) {
	// 避免级联错误.
	if c.panicMode {
		return
	}
	c.panicMode = true
	fmt.Printf("[line %d] Error", token.line)
	if token.typ == TOKEN_EOF {
		fmt.Printf(" at end")
	} else if token.typ == TOKEN_ERROR {
		// Nothing.
	} else {
		fmt.Printf(" at '%s'", token.value)
	}
	fmt.Printf(": %s\n", message)
	c.hadError = true
}

// 对每个token，作为前缀表达式的函数和中缀表达式的函数.
func (c *Compiler) createRules() []*ParseRule {
	return []*ParseRule{
		NewParseRule(c.grouping, nil, PREC_NONE),     // TOKEN_LEFT_PAREN
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_RIGHT_PAREN
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_LEFT_BRACE
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_RIGHT_BRACE
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_COMMA
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_DOT
		NewParseRule(c.unary, c.binary, PREC_TERM),   // TOKEN_MINUS
		NewParseRule(nil, c.binary, PREC_TERM),       // TOKEN_PLUS
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_SEMICOLON
		NewParseRule(nil, c.binary, PREC_FACTOR),     // TOKEN_SLASH
		NewParseRule(nil, c.binary, PREC_FACTOR),     // TOKEN_STAR
		NewParseRule(c.unary, nil, PREC_NONE),        // TOKEN_BANG
		NewParseRule(nil, c.binary, PREC_EQUALITY),   // TOKEN_BANG_EQUAL
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_EQUAL
		NewParseRule(nil, c.binary, PREC_EQUALITY),   // TOKEN_EQUAL_EQUAL
		NewParseRule(nil, c.binary, PREC_COMPARISON), // TOKEN_GREATER
		NewParseRule(nil, c.binary, PREC_COMPARISON), // TOKEN_GREATER_EQUAL
		NewParseRule(nil, c.binary, PREC_COMPARISON), // TOKEN_LESS
		NewParseRule(nil, c.binary, PREC_COMPARISON), // TOKEN_LESS_EQUAL
		NewParseRule(c.variable, nil, PREC_NONE),     // TOKEN_IDENTIFIER
		NewParseRule(c.string, nil, PREC_NONE),       // TOKEN_STRING
		NewParseRule(c.number, nil, PREC_NONE),       // TOKEN_NUMBER
		NewParseRule(nil, c.and_, PREC_AND),          // TOKEN_AND
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_CLASS
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_ELSE
		NewParseRule(c.literal, nil, PREC_NONE),      // TOKEN_FALSE
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_FOR
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_FUN
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_IF
		NewParseRule(c.literal, nil, PREC_NONE),      // TOKEN_NIL
		NewParseRule(nil, c.or_, PREC_OR),            // TOKEN_OR
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_PRINT
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_RETURN
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_SUPER
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_THIS
		NewParseRule(c.literal, nil, PREC_NONE),      // TOKEN_TRUE
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_VAR
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_WHILE
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_ERROR
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_EOF
	}
}

func (c *Compiler) getRule(typ TokenType) *ParseRule {
	return c.rules[typ]
}

type state struct {
	hadError bool
	// 当前是否在紧急模式中，即跳过错误的代码直到遇到下一条`语句`
	// 到达一个同步点时，紧急模式就结束了
	panicMode bool
	previous  *Token
	current   *Token
}

func NewState() *state {
	return &state{}
}

type Precedence byte

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . () []
	PREC_PRIMARY
)

type Parselet func(canAssign bool)
type ParseRule struct {
	prefix     Parselet   // 以该类型标识为起点的前缀表达式的函数
	infix      Parselet   // 左操作数后跟该类型标识的中缀表达式的函数
	precedence Precedence // 使用该标识作为操作符的`中缀表达式`的优先级
}

func NewParseRule(prefix Parselet, infix Parselet, precedence Precedence) *ParseRule {
	return &ParseRule{
		prefix:     prefix,
		infix:      infix,
		precedence: precedence,
	}
}

type Resolver struct {
	function *ObjFunction
	typ      FunctionType

	locals     [UINT8_COUNT]*Local // 越往后，作用域深度越大
	localCount int                 // 当前作用域的局部变量数
	scopeDepth int                 // 当前作用域的深度
}

func NewResolver(typ FunctionType) *Resolver {
	res := &Resolver{function: NewObjFunction("", 0, NewChunk()), typ: typ}
	for i := 0; i < len(res.locals); i++ {
		res.locals[i] = NewLocal()
	}
	res.localCount++ // 将栈槽0用于虚拟机的内部使用
	return res
}

type Local struct {
	name  *Token
	depth int
}

func NewLocal() *Local {
	return &Local{}
}

func (l *Local) String() string {
	return fmt.Sprintf("%s", l.name.value)
}

// #endregion

// #region Scanner 扫描器生成tokens
type Scanner struct {
	start   int // start of lexeme
	current int // current character of lexeme
	line    int // current line of lexeme

	source string
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source}
}

// 用于测试.
func (s *Scanner) Scan() {
	line := -1
	for {
		token := s.ScanToken()
		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%s'\n", token.typ, token.value)

		if token.typ == TOKEN_EOF {
			break
		}
	}
}

// 每次调用都会扫描一个完整的词法标识，所以当我们进入该函数时，就知道我们正处于一个新词法标识的开始处.
func (s *Scanner) ScanToken() *Token {
	s.skipWhitespace()
	s.start = s.current
	if s.isAtEnd() {
		return s.makeToken(TOKEN_EOF)
	}

	c := s.advance()
	if isAlpha(c) {
		return s.identifier()
	}
	if isDigit(c) {
		return s.number()
	}

	switch c {
	case '(':
		return s.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return s.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return s.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return s.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return s.makeToken(TOKEN_SEMICOLON)
	case ',':
		return s.makeToken(TOKEN_COMMA)
	case '.':
		return s.makeToken(TOKEN_DOT)
	case '-':
		return s.makeToken(TOKEN_MINUS)
	case '+':
		return s.makeToken(TOKEN_PLUS)
	case '/':
		return s.makeToken(TOKEN_SLASH)
	case '*':
		return s.makeToken(TOKEN_STAR)

		// !多字符标记
	case '!':
		if s.match('=') {
			return s.makeToken(TOKEN_BANG_EQUAL)
		} else {
			return s.makeToken(TOKEN_BANG)
		}
	case '=':
		if s.match('=') {
			return s.makeToken(TOKEN_EQUAL_EQUAL)
		} else {
			return s.makeToken(TOKEN_EQUAL)
		}
	case '<':
		if s.match('=') {
			return s.makeToken(TOKEN_LESS_EQUAL)
		} else {
			return s.makeToken(TOKEN_LESS)
		}
	case '>':
		if s.match('=') {
			return s.makeToken(TOKEN_GREATER_EQUAL)
		} else {
			return s.makeToken(TOKEN_GREATER)
		}

		// !字符串
	case '"':
		return s.string()
	}

	return s.errorToken("Unexpected character.")
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current+1]
}

// 如果当前字符是所需的字符，则指针前进并返回true。否则，我们返回false表示没有匹配。
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) isAtEnd() bool {
	return s.current == len(s.source)
}

func (s *Scanner) makeToken(typ TokenType) *Token {
	return &Token{
		typ:   typ,
		value: s.source[s.start:s.current],
		line:  s.line,
	}
}

func (s *Scanner) errorToken(message string) *Token {
	return &Token{
		typ:   TOKEN_ERROR,
		value: message,
		line:  s.line,
	}
}

func (s *Scanner) skipWhitespace() {
	for {
		c := s.peek()
		switch c {
		case ' ', '\r', '\t':
			s.advance()
			break
		case '\n':
			s.line++
			s.advance()
			break
		case '/':
			// A comment goes until the end of the line.
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
			break
		default:
			return
		}
	}
}

func (s *Scanner) identifier() *Token {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}
	return s.makeToken(s.identifierType())
}

// 简化版trie.
func (s *Scanner) identifierType() TokenType {
	switch s.source[s.start] {
	case 'a':
		return s.checkKeyWord(1, 2, "nd", TOKEN_AND)
	case 'c':
		return s.checkKeyWord(1, 4, "lass", TOKEN_CLASS)
	case 'e':
		return s.checkKeyWord(1, 3, "lse", TOKEN_ELSE)
	case 'f':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyWord(2, 3, "lse", TOKEN_FALSE)
			case 'o':
				return s.checkKeyWord(2, 1, "r", TOKEN_FOR)
			case 'u':
				return s.checkKeyWord(2, 1, "n", TOKEN_FUN)
			}
		}
		break
	case 'i':
		return s.checkKeyWord(1, 1, "f", TOKEN_IF)
	case 'n':
		return s.checkKeyWord(1, 2, "il", TOKEN_NIL)
	case 'o':
		return s.checkKeyWord(1, 1, "r", TOKEN_OR)
	case 'p':
		return s.checkKeyWord(1, 4, "rint", TOKEN_PRINT)
	case 'r':
		return s.checkKeyWord(1, 5, "eturn", TOKEN_RETURN)
	case 's':
		return s.checkKeyWord(1, 4, "uper", TOKEN_SUPER)
	case 't':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyWord(2, 2, "is", TOKEN_THIS)
			case 'r':
				return s.checkKeyWord(2, 2, "ue", TOKEN_TRUE)
			}
		}
		break
	case 'v':
		return s.checkKeyWord(1, 2, "ar", TOKEN_VAR)
	case 'w':
		return s.checkKeyWord(1, 4, "hile", TOKEN_WHILE)
	}
	return TOKEN_IDENTIFIER
}

func (s *Scanner) checkKeyWord(offset, length int, rest string, typ TokenType) TokenType {
	if s.current-s.start == offset+length && s.source[s.start+offset:s.start+offset+length] == rest {
		return typ
	}
	return TOKEN_IDENTIFIER
}

func (s *Scanner) number() *Token {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	return s.makeToken(TOKEN_NUMBER)
}

func (s *Scanner) string() *Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}
	s.advance() // closing quote
	return s.makeToken(TOKEN_STRING)
}

// #endregion

// #region Token
type Token struct {
	typ   TokenType
	value string
	line  int
}

type TokenType byte

const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUN
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_ERROR
	TOKEN_EOF
)

// #endregion

// #region Utils
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func hashString(s string) int {
	res := 0
	for i := 0; i < len(s); i++ {
		res = res*31 + int(s[i])
	}
	return res
}

// #endregion
