/**
 * 词法解析器
 * 在 01.js 的基础上，支持 LL(k), k > 1。递归下降解析
 *
 */
const EOF = -1
const EOF_TYPE = 1

class Lexer {
	constructor(input) {
		this.input = input // 输入的字符串
		this.index = 0 // 当前字符串的索引位置
		this.char = input[this.index] // 当前字符
	}
	consume() { // 向前移动一个字符
		this.index += 1
		if (this.index >= this.input.length) { // 判断是否到末尾
			this.char = EOF
		} else {
			this.char = this.input[this.index]
		}
	}
	match(char) { // 判断输入的 char 是否为当前的 this.char
		if (this.char === char) {
			this.consume()
		} else {
			throw new Error(`Expecting ${char}; Found ${this.char}`)
		}
	}
}

Lexer.EOF = EOF
Lexer.EOF_TYPE = EOF_TYPE

const NAME = 2
const COMMA = 3
const LBRACK = 4
const RBRACK = 5
const EQUALS = 6
const tokenNames = ['n/a', '<EOF>', 'NAME', 'COMMA', 'LBRACK', 'RBRACK', 'EQUALS']
const getTokenName = index => tokenNames[index]

// 判断输入字符是否为字母，即在 a-zA-Z 之间
const isLetter = char => char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'

class ListLexer extends Lexer {
	constructor(input) {
		super(input)
	}
	isLetter() {
		return isLetter(this.char)
	}
	nextToken() {
		while (this.char !== EOF) {
			switch (this.char) {
				case ' ':
				case '\t':
				case '\n':
				case '\r':
					this.WS()
					break
				case ',':
					this.consume()
					return new Token(COMMA, ',')
				case '[':
					this.consume()
					return new Token(LBRACK, '[')
				case ']':
					this.consume()
					return new Token(RBRACK, ']')
				case '=':
					this.consume()
					return new Token(EQUALS, '=')
				default:
					if (this.isLetter()) {
						return this.NAME()
					}
					throw new Error(`Invalid character: ${this.char}`)
			}
		}
		return new Token(EOF_TYPE, '<EOF>')
	}
	WS() { // 忽略所有空白，换行，tab，回车符等
		while (this.char === ' ' || this.char === '\t' || this.char === '\n' || this.char === '\r') {
			this.consume()
		}
	}
	NAME() { // 匹配一列字母
		let name = ''
		while (this.isLetter()) {
			name += this.char
			this.consume()
		}
		return new Token(NAME, name)
	}
}

class Token {
	constructor(type, text) {
		this.type = type
		this.text = text
	}
	toString() {
		let tokenName = tokenNames[this.type]
		return `<'${this.text}',${tokenName}>`
	}
}


class Parser {
	constructor(lexer, k) { // lexer 词法解析类的实例，k 向前预读的 token 数量
		this.lexer = lexer
		this.k = k
		this.index = 0
		this.lookahead = Array.from(Array(k)) // 预读 token 列表
		this.lookahead.forEach(() => { // 调用 consume 方法填充列表
			this.consume()
		})
	}
	consume() {
		this.lookahead[this.index] = this.lexer.nextToken()
		this.index = (this.index + 1) % this.k
	}
	getToken(n) {
		let index = (this.index + n - 1) % this.k
		return this.lookahead[index]
	}
	getTokenType(n) {
		return this.getToken(n).type
	}
	match(type) {
		let tokenType = this.getTokenType(1)
		if (tokenType === type) {
			this.consume()
		} else {
			throw new Error(`Expecting ${getTokenName(type)}; Found ${this.getToken()}`)
		}
	}
}


class ListParser extends Parser {
	list() {
		this.match(LBRACK)
		this.elements()
		this.match(RBRACK)
	}
	elements() {
		this.element()
		while (this.getTokenType(1) === COMMA) {
			this.match(COMMA)
			this.element()
		}
	}
	element() {
		let tokenType1 = this.getTokenType(1)
		let tokenType2 = this.getTokenType(2)

		if (tokenType1 === NAME && tokenType2 === EQUALS) {
			this.match(NAME)
			this.match(EQUALS)
			this.match(NAME)
		} else if (tokenType1 === NAME) {
			this.match(NAME)
		} else if (tokenType1 === LBRACK) {
			this.list()
		} else {
			throw new Error(`Expecting name or list; Found ${this.getToken(1)}`)
		}
	}
}




exports.ListLexer = ListLexer
exports.ListParser = ListParser