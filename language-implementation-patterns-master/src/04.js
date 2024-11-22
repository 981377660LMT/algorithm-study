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

// 是数字字符 0~9
const isNumber = char => char >= '0' && char <= '9'

// 检查 input 是否等于 char，value 为之前检查通过的 char 拼接起来的字符串
const createMatch = (char, length) => (input, value) => {
	if (typeof length === 'number' && value.length === length) {
		return false
	}
	return input === char
}

const createFlatTokenizer = match => (input, start) => {
	if (!match.start) {
		let _match = match
		match = {
			start: _match,
			end() {
				return !_match.apply(null, arguments)
			}
		}
	}

	let offset = 0
	let item = input[start + offset]
	let value = []

	if (!match.start(item, value)) {
		return
	}

	do {
		value.push(item)
		offset += 1
		item = input[start + offset]
	} while (!match.end(item, value))

	return {
		value: value,
		start: start,
		end: start + offset
	}
}

const firstLayerPatterns = {
	'LETTER': createFlatTokenizer(isLetter),
	'WHITE_SPACE': createFlatTokenizer(isWhiteSpace),
	'NUMBER': createFlatTokenizer(isNumber),
	'DOUBLE_QUOTES': createFlatTokenizer(createMatch('"', 1)),
	'SINGLE_QUOTE': createFlatTokenizer(createMatch('\'', 1)),
	'LEFT_BRACKET': createFlatTokenizer(createMatch('(', 1)),
	'RIGHT_BRACKET': createFlatTokenizer(createMatch(')', 1)),
	'LEFT_BRACE': createFlatTokenizer(createMatch('{', 1)),
	'RIGHT_BRACE': createFlatTokenizer(createMatch('}', 1)),
	'SLASH': createFlatTokenizer(createMatch('/', 1)),
	'BACK_SLANT': createFlatTokenizer(createMatch('\\', 1)),
	'WHIFFLETREE': createFlatTokenizer(createMatch('-', 1)),
	'BANG': createFlatTokenizer(createMatch('!', 1)),
	'COLON': createFlatTokenizer(createMatch(':', 1)),
	'DOT': createFlatTokenizer(createMatch('.', 1)),
	'POUND_KEY': createFlatTokenizer(createMatch('#', 1)),
	'SEMICOLON': createFlatTokenizer(createMatch(';', 1)),
	'AT_SYMBOL': createFlatTokenizer(createMatch('@', 1)),
	'COMMA': createFlatTokenizer(createMatch(',', 1)),
	'EQUAL_SYMBOL': createFlatTokenizer(createMatch('=', 1)),
	'UNDERLINE': createFlatTokenizer(createMatch('_', 1)),
	'PERSCENT_SYMBOL': createFlatTokenizer(createMatch('%', 1)),
	'ASTERISK': createFlatTokenizer(createMatch('*', 1)),
	'LEFT_ANGLE_BRACKET': createFlatTokenizer(createMatch('<', 1)),
	'RIGHT_ANGLE_BRACKET': createFlatTokenizer(createMatch('>', 1)),
}

function getFirstLayerToken(input, start) {
	for (let key in firstLayerPatterns) {
		let tokenizer = firstLayerPatterns[key]
		let token = tokenizer(input, start)
		if (token) {
			return Object.assign({type: key}, token)
		}
	}
}

function firstLayerTokenizer(input) {
	let index = 0
	let tokens = []

	while (index < input.length) {
		let token = getFirstLayerToken(input, index)

		if (!token) {
			throw new Error(`Unknow char: ${input[index]}`)
		}

		tokens.push(token)

		index = token.end
	}

	return tokens
}

let cssFilePath = path.join(__dirname, 'files/test.css')
let content = fs.readFileSync(cssFilePath).toString()
let tokens = firstLayerTokenizer(content)

const destPath = path.join(__dirname, 'dest/04.txt')
fs.writeFileSync(destPath, JSON.stringify(tokens, null, 2))