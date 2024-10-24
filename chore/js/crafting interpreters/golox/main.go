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
	Run("!(5 - 4 > 3 * 2 == !nil)")

}

func TestScanner(source string) {
	S := NewScanner(source)
	S.Scan()
}

func Run(source string) {
	scanner := NewScanner(source)
	compiler := NewCompiler(scanner)

	vm.Init()
	vm.Inject(compiler)
	res := vm.Interpret(true)
	fmt.Println(res)
}

// #region Chunk

type OpCode = byte // type of bytecode instructions

const (
	OP_CONSTANT OpCode = iota // load the constant for use
	OP_NIL
	OP_TRUE
	OP_FALSE
	OP_EQUAL
	OP_GREATER
	OP_LESS
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NOT
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

func (c *Chunk) Write(b byte, line int) {
	c.code = append(c.code, b)
	c.lines = append(c.lines, line)
}

func (c *Chunk) AddConstant(v float64) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
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
	case OP_NIL:
		return c.simpleInstruction("OP_NIL", offset)
	case OP_TRUE:
		return c.simpleInstruction("OP_TRUE", offset)
	case OP_FALSE:
		return c.simpleInstruction("OP_FALSE", offset)
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

// #region Value
type ValueType byte

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
)

type Value struct {
	typ ValueType
	// union
	asBool   bool
	asNumber float64
}

func NewValueNumber(v float64) Value {
	return Value{typ: VAL_NUMBER, asNumber: v}
}

func NewValueBool(v bool) Value {
	return Value{typ: VAL_BOOL, asBool: v}
}

func NewValueNil() Value {
	return Value{typ: VAL_NIL}
}

func IsSameValue(a, b Value) bool {
	if a.typ != b.typ {
		return false
	}
	switch a.typ {
	case VAL_BOOL:
		return a.asBool == b.asBool
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return a.asNumber == b.asNumber
	}
	return false
}

func IsBool(v Value) bool {
	return v.typ == VAL_BOOL
}

func IsNil(v Value) bool {
	return v.typ == VAL_NIL
}

func IsNumber(v Value) bool {
	return v.typ == VAL_NUMBER
}

func AsBool(v Value) bool {
	return v.asBool
}

func AsNumber(v Value) float64 {
	return v.asNumber
}

func (v Value) String() string {
	switch v.typ {
	case VAL_BOOL:
		if v.asBool {
			return "true"
		} else {
			return "false"
		}
	case VAL_NIL:
		return "nil"
	case VAL_NUMBER:
		return fmt.Sprintf("%g", v.asNumber)
	}
	return ""
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
	// chunk.code[ip]
	ip int

	stack    [STACK_MAX]Value
	stackTop int // points to where the next value to be pushed will go

	compiler *Compiler
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) Init() {
}

func (vm *VM) Inject(compiler *Compiler) {
	vm.compiler = compiler
}

func (vm *VM) Interpret(debug bool) InterpretResult {
	chunk := NewChunk()
	if !vm.compiler.Compile(chunk) {
		return INTERPRET_COMPILE_ERROR
	}
	vm.chunk = chunk
	vm.ip = 0
	res := vm.run(debug)
	return res
}

func (vm *VM) run(debug bool) InterpretResult {
	if debug {
		fmt.Println("         ")
		for i := 0; i < vm.stackTop; i++ {
			fmt.Printf("[ %s ]\n", vm.stack[i])
		}
		vm.chunk.Disassemble("test chunk")
	}

	for {
		// decoding/dispatching
		instruction := vm.readByte()
		switch instruction {
		case OP_CONSTANT:
			constant := vm.readConstant()
			vm.push(NewValueNumber(constant))
			break
		case OP_NIL:
			vm.push(NewValueNil())
			break
		case OP_TRUE:
			vm.push(NewValueBool(true))
			break
		case OP_FALSE:
			vm.push(NewValueBool(false))
			break
		case OP_EQUAL:
			b := vm.pop()
			a := vm.pop()
			vm.push(NewValueBool(IsSameValue(a, b)))
			break
		case OP_GREATER:
			vm.binaryOp(func(a, b float64) Value { return NewValueBool(a > b) })
			break
		case OP_LESS:
			vm.binaryOp(func(a, b float64) Value { return NewValueBool(a < b) })
			break
		case OP_ADD:
			vm.binaryOp(func(a, b float64) Value { return NewValueNumber(a + b) })
			break
		case OP_SUBTRACT:
			vm.binaryOp(func(a, b float64) Value { return NewValueNumber(a - b) })
			break
		case OP_MULTIPLY:
			vm.binaryOp(func(a, b float64) Value { return NewValueNumber(a * b) })
			break
		case OP_DIVIDE:
			vm.binaryOp(func(a, b float64) Value { return NewValueNumber(a / b) })
			break
		case OP_NOT:
			vm.push(NewValueBool(vm.isFalsey(vm.pop())))
			break
		case OP_NEGATE:
			if !IsNumber(vm.peek(0)) {
				vm.runtimeError("Operand must be a number.")
				return INTERPRET_RUNTIME_ERROR
			}
			vm.push(NewValueNumber(-AsNumber(vm.pop())))
			break
		case OP_RETURN:
			fmt.Println(vm.pop())
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

func (vm *VM) push(v Value) {
	vm.stack[vm.stackTop] = v
	vm.stackTop++
}

func (vm *VM) pop() Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

// 将操作数留在栈上是很重要的，可以确保在运行过程中触发垃圾收集时，垃圾收集器能够找到它们.
func (vm *VM) peek(distance int) Value {
	return vm.stack[vm.stackTop-1-distance]
}

func (vm *VM) isFalsey(v Value) bool {
	return IsNil(v) || (IsBool(v) && !AsBool(v))
}

func (vm *VM) binaryOp(f func(float64, float64) Value) {
	if !IsNumber(vm.peek(0)) || !IsNumber(vm.peek(1)) {
		vm.runtimeError("Operands must be numbers.")
		return
	}
	b := AsNumber(vm.pop())
	a := AsNumber(vm.pop())
	vm.push(f(a, b))
}

func (vm *VM) runtimeError(message string) {
	fmt.Printf("%s\n", message)

	instruction := vm.ip
	line := vm.chunk.lines[instruction]
	fmt.Printf("[line %d] in script\n", line)
}

// #endregion

// #region Compiler
type Compiler struct {
	*state
	rules          []*ParseRule
	compilingChunk *Chunk

	scanner *Scanner
}

func NewCompiler(scanner *Scanner) *Compiler {
	res := &Compiler{state: NewState(), scanner: scanner}
	res.rules = res.createRules()
	return res
}

// 解析源代码并输出低级的二进制指令序列。
// 当然，它是字节码，而不是某个芯片的原生指令集，但它比jlox更接近于硬件。
// !我们将字节码块传入，编译器会向其中写入代码。返回是否编译成功。
func (c *Compiler) Compile(chunk *Chunk) bool {
	c.compilingChunk = chunk
	c.advance()
	c.expression()
	// 在编译表达式之后，我们应该处于源代码的末尾，所以我们要检查EOF标识。
	c.consume(TOKEN_EOF, "Expect end of expression.")
	c.endCompiler()
	return !c.hadError
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

// 读取下一个标识、验证标识是否具有预期的类型。如果不是，则报告错误。
func (c *Compiler) consume(tokenType TokenType, message string) {
	if c.current.typ == tokenType {
		c.advance()
		return
	}
	c.errorAtCurrent(message)
}

// 向块中追加一个字节。
// 将给定的字节写入一个指令，该字节可以是操作码或操作数。
// 它会发送前一个token的行信息，以便将运行时错误与该行关联起来。
func (c *Compiler) emitByte(b byte) {
	c.currentChunk().Write(b, c.previous.line)
}

// 一般是写一个操作码，后面跟一个单字节的操作数。
func (c *Compiler) emitByte2(b1, b2 byte) {
	c.emitByte(b1)
	c.emitByte(b2)
}

func (c *Compiler) emitReturn() {
	c.emitByte(OP_RETURN)
}

func (c *Compiler) makeConstant(value Value) byte {
	constant := c.currentChunk().AddConstant(AsNumber(value))
	// OP_CONSTANT指令使用单个字节来索引操作数，所以我们在一个块中最多只能存储和加载256个常量。
	if constant > math.MaxUint8 {
		c.error("Too many constants in one chunk")
		return 0
	}
	return byte(constant)
}

func (c *Compiler) emitConstant(value Value) {
	c.emitByte2(OP_CONSTANT, c.makeConstant(value))
}

func (c *Compiler) endCompiler() {
	c.emitReturn()
}

func (c *Compiler) binary() {
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

func (c *Compiler) literal() {
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

func (c *Compiler) grouping() {
	// 假定初始的(已经被消耗了。递归地调用expression()来编译括号之间的表达式，然后解析结尾的)。
	c.expression()
	c.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (c *Compiler) number() {
	num, _ := strconv.ParseFloat(c.previous.value, 64)
	c.emitConstant(NewValueNumber(num))
}

func (c *Compiler) unary() {
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
	prefixParselet()
	for p <= c.getRule(c.current.typ).precedence {
		c.advance()
		infixParselet := c.getRule(c.previous.typ).infix
		infixParselet()
	}
}

func (c *Compiler) currentChunk() *Chunk {
	return c.compilingChunk
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
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_IDENTIFIER
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_STRING
		NewParseRule(c.number, nil, PREC_NONE),       // TOKEN_NUMBER
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_AND
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_CLASS
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_ELSE
		NewParseRule(c.literal, nil, PREC_NONE),      // TOKEN_FALSE
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_FOR
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_FUN
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_IF
		NewParseRule(c.literal, nil, PREC_NONE),      // TOKEN_NIL
		NewParseRule(nil, nil, PREC_NONE),            // TOKEN_OR
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

type Parselet func()
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

// #endregion

// #region Scanner
type Scanner struct {
	start   int // start of lexeme
	current int // current character of lexeme
	line    int // current line of lexeme

	source string
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source}
}

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

// #endregion
