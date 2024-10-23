package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	Run("3 + 4 / 5") // expect 23
	// TestScanner("print 1 + 2;")
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
	// chunk.code[ip]
	ip int

	stack    [STACK_MAX]float64
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
		if c.current.kind != TOKEN_ERROR {
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
	if c.current.kind == tokenType {
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

func (c *Compiler) makeConstant(v float64) byte {
	constant := c.currentChunk().AddConstant(v)
	// OP_CONSTANT指令使用单个字节来索引操作数，所以我们在一个块中最多只能存储和加载256个常量。
	if constant > math.MaxUint8 {
		c.error("Too many constants in one chunk")
		return 0
	}
	return byte(constant)
}

func (c *Compiler) emitConstant(v float64) {
	c.emitByte(OP_RETURN)
}

func (c *Compiler) endCompiler() {
	c.emitReturn()
}

func (c *Compiler) binary() {
	kind := c.previous.kind
	// 为每个二元运算符定义一个单独的函数。每个函数都会调用 parsePrecedence() 并传入正确的优先级来解析其操作数。
	rule := c.getRule(kind)
	c.parsePrecedence(rule.precedence + 1)

	switch kind {
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

func (c *Compiler) grouping() {
	// 假定初始的(已经被消耗了。递归地调用expression()来编译括号之间的表达式，然后解析结尾的)。
	c.expression()
	c.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (c *Compiler) number() {
	num, _ := strconv.ParseFloat(c.previous.value, 64)
	c.emitByte2(OP_CONSTANT, c.makeConstant(num))
}

func (c *Compiler) unary() {
	kind := c.previous.kind
	c.parsePrecedence(PREC_UNARY)
	switch kind {
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
	prefixParselet := c.getRule(c.previous.kind).prefix
	if prefixParselet == nil {
		c.error("Expect expression.")
		return
	}
	prefixParselet()
	for p <= c.getRule(c.current.kind).precedence {
		c.advance()
		infixParselet := c.getRule(c.previous.kind).infix
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
	if token.kind == TOKEN_EOF {
		fmt.Printf(" at end")
	} else if token.kind == TOKEN_ERROR {
		// Nothing.
	} else {
		fmt.Printf(" at '%s'", token.value)
	}
	fmt.Printf(": %s\n", message)
	c.hadError = true
}

func (c *Compiler) createRules() []*ParseRule {
	return []*ParseRule{
		NewParseRule(c.grouping, nil, PREC_NONE),   // TOKEN_LEFT_PAREN
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_RIGHT_PAREN
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_LEFT_BRACE
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_RIGHT_BRACE
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_COMMA
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_DOT
		NewParseRule(c.unary, c.binary, PREC_TERM), // TOKEN_MINUS
		NewParseRule(nil, c.binary, PREC_TERM),     // TOKEN_PLUS
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_SEMICOLON
		NewParseRule(nil, c.binary, PREC_FACTOR),   // TOKEN_SLASH
		NewParseRule(nil, c.binary, PREC_FACTOR),   // TOKEN_STAR
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_BANG
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_BANG_EQUAL
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_EQUAL
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_EQUAL_EQUAL
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_GREATER
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_GREATER_EQUAL
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_LESS
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_LESS_EQUAL
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_IDENTIFIER
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_STRING
		NewParseRule(c.number, nil, PREC_NONE),     // TOKEN_NUMBER
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_AND
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_CLASS
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_ELSE
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_FALSE
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_FOR
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_FUN
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_IF
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_N
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_OR
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_PRINT
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_RETURN
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_SUPER
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_THIS
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_TRUE
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_VAR
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_WHILE
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_ERROR
		NewParseRule(nil, nil, PREC_NONE),          // TOKEN_EOF
	}
}

func (c *Compiler) getRule(kind TokenType) *ParseRule {
	return c.rules[kind]
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
		fmt.Printf("%2d '%s'\n", token.kind, token.value)

		if token.kind == TOKEN_EOF {
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

func (s *Scanner) makeToken(kind TokenType) *Token {
	return &Token{
		kind:  kind,
		value: s.source[s.start:s.current],
		line:  s.line,
	}
}

func (s *Scanner) errorToken(message string) *Token {
	return &Token{
		kind:  TOKEN_ERROR,
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

func (s *Scanner) checkKeyWord(offset, length int, rest string, kind TokenType) TokenType {
	if s.current-s.start == offset+length && s.source[s.start+offset:s.start+offset+length] == rest {
		return kind
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
	kind  TokenType
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
