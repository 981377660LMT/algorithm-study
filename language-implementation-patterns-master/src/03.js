const fs = require('fs')
const path = require('path')

// 组合两个函数
const compose = (f1, f2) => arg => f1(f2(arg))

// 连接两个函数
const pipe = (f1, f2) => arg => f2(f1(arg))

const isFn = obj => typeof obj === 'function'

const returnTrue = () => true

// 判断输入字符是否为字母，即在 a-zA-Z 之间
const isLetter = char => char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'

// 是空格或者其他分割符
const isWhiteSpace = char => char === ' ' || char === '\t' || char === '\n' || char === '\r'

const isNumber = char => !isNaN(Number(char.trim()))

const charMap = {
	'LETTER': isLetter,
	'WHITE_SPACE': isWhiteSpace,
	'NUMBER': isNumber,
	'DOUBLE_QUOTES': '"',
	'SINGLE_QUOTE': '\'',
	'LEFT_BRACKET': '(',
	'RIGHT_BRACKET': ')',
	'LEFT_BRACE': '{',
	'RIGHT_BRACE': '}',
	'SLASH': '/',
	'BACK_SLANT': '\\',
	'WHIFFLETREE': '-',
	'BANG': '!',
	'COLON': ':',
	'DOT': '.',
	'POUND_KEY': '#',
	'SEMICOLON': ';',
	'AT_SYMBOL': '@',
	'COMMA': ',',
	'EQUAL_SYMBOL': '=',
	'UNDERLINE': '_',
	'PERSCENT_SYMBOL': '%',
	'ASTERISK': '*',
}

const getCharPattern = char => {
	for (let key in charMap) {
		let value = charMap[key]
		if (isFn(value) && value(char)) {
			return {
				type: key,
				match: value,
			}
		} else if (char === value) {
			return {
				type: key,
				match: char => char === value,
			}
		}
	}

	// default
	return {
		type: 'UNKNOW',
		match: returnTrue
	}
}

class CharStream {
	constructor(input) {
		this.input = input
		this.index = 0
	}
	current() {
		return this.input[this.index]
	}
	next() {
		if (this.isEnd()) {
			return
		}
		this.index += 1
		return this.current()
	}
	test(char) {
		return char === this.current()
	}
	isEnd() {
		return this.index >= this.input.length
	}
}

class Token {
	constructor(type, text) {
		this.type = type
		this.text = text
	}
	toString() {
		return `${this.type}: ${this.text}`
	}
	length() {
		return this.text.length
	}
}

class TokenStream {
	constructor(input) {
		this.charStream = new CharStream(input)
		this.index = 0
		this.token = null
	}
	isEnd() {
		return this.charStream.isEnd()
	}
	current() {
		return this.token
	}
	next() {
		if (this.isEnd()) {
			return
		}
		let char = this.charStream.current()
		let { type, match } = getCharPattern(char)
		this.token = this.createToken(type, match)
		return this.current()
	}
	createToken(type, match) {
		let { charStream } = this
		let text = ''
		let char = charStream.current()
		while (!charStream.isEnd() && match(char)) {
			text += char
			char = charStream.next()
		}
		return new Token(type, text)
	}
}


// let cssFilePath = path.join(__dirname, 'files/test.css')
// let content = fs.readFileSync(cssFilePath).toString()
// let tokenStream = new TokenStream(content)

// let result = ''
// while (!tokenStream.isEnd()) {
// 	let token = tokenStream.next().toString()
// 	result += token + '\n'
// }

// const destPath = path.join(__dirname, 'dest/03.txt')
// fs.writeFileSync(destPath, result)


// const isClassSelector = 

// const syntaxMap = {
// 	'SELECTOR': `
// 		[DOT|POUND_KEY]LETTER
// 	`,
// }

// class SyntaxStream {
// 	constructor(input) {
// 		this.tokenStream = new TokenStream(input)
// 		this.tokens = []
// 		this.index = 0
// 	}
// 	isEnd() {
// 		return this.tokenStream.isEnd()
// 	}
// 	getCurrentToken() {
// 		return this.tokens[this.index]
// 	}
// 	getNextToken() {
// 		if (this.isEnd()) {
// 			return
// 		} 
// 		let token = this.tokenStream.next()
// 		this.tokens.push(token)
// 		return token
// 	}
// 	isClassSelector() {
// 		let offset = 0
// 		let currentToken = this.getCurrentToken()
// 		if (currentToken.type === 'DOT') {
			
// 		}
// 	}
// }

