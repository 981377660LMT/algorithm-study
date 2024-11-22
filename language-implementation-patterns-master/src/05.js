// 判断输入字符是否为字母，即在 a-zA-Z 之间
const isLetter = char => char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'

// 是空格或者其他分割符
const isWhiteSpace = char => char === ' ' || char === '\t' || char === '\n' || char === '\r'

// 是数字字符 0~9
const isNumber = char => char >= '0' && char <= '9'

// 检查 input 是否等于 char，value 为之前检查通过的 char 拼接起来的字符串
const createMatch = (char, length) => {
	let start = value => value === char
	let end = value => {
		if (typeof length === 'number' && value.length === length) {
			return true
		}
		return value !== char
	}
	return { start, end }
}

const defaultGetValue = (item, value) => value + item

const createTokenizer = match => {
	let tokenizer = (input, start) => {
		let item = input[start]
		// if (match.test) {
		// 	console.log(item, start)
		// }
		if (!match.start(item)) {
			return
		}

		let value = ''

		if (match.getInitialValue) {
			value = match.getInitialValue(item)
		}

		let getValue = match.getValue || defaultGetValue
		let offset = 0

		if (match.first !== false) {
			value = getValue(item, value, offset)
		}

		if (match.check && match.check(value, offset, start) === false) {
			return
		}

		offset += 1
		item = input[start + offset]

		while (item !== undefined && !match.end(item, value)) {
			value = getValue(item, value, offset)
			offset += 1
			item = input[start + offset]
			if (match.check && match.check(value, offset, start) === false) {
				return
			}
		}

		if (item && match.last === true) {
			value = getValue(item, value, offset)
			offset += 1
			if (match.check && match.check(value, offset, start) === false) {
				return
			}
		}

		if (match.getFinalValue) {
			value = match.getFinalValue(value, offset, start)
		}

		return {
			value: value,
			start: start,
			end: start + offset
		}
	}

	return tokenizer
}

const patterns = {
	'DOUBLE_QUOTES_STRING': createTokenizer({
		last: true,
		start: char => char === '"',
		end: char => char === '"',
	}),
	'SINGLE_QUOTES_STRING': createTokenizer({
		last: true,
		start: char => char === '\'',
		end: char => char === '\'',
	}),
	'NAME': createTokenizer({
		start: isLetter,
		end: char => 
			!isLetter(char) &&
			!isNumber(char) &&
			char !== '_' &&
			char !== '-',
	}),
	'NUMBER': createTokenizer({
		start: isNumber,
		end: value => !isNumber(value),
	}),
	'WHITE_SPACE': createTokenizer({
		start: isWhiteSpace,
		end: value => !isWhiteSpace(value),
	}),
	'LEFT_BRACKET': createTokenizer(createMatch('(', 1)),
	'RIGHT_BRACKET': createTokenizer(createMatch(')', 1)),
	'LEFT_BRACE': createTokenizer(createMatch('{', 1)),
	'RIGHT_BRACE': createTokenizer(createMatch('}', 1)),
	'SLASH': createTokenizer(createMatch('/', 1)),
	'BACK_SLANT': createTokenizer(createMatch('\\', 1)),
	'WHIFFLETREE': createTokenizer(createMatch('-', 1)),
	'BANG': createTokenizer(createMatch('!', 1)),
	'COLON': createTokenizer(createMatch(':', 1)),
	'DOT': createTokenizer(createMatch('.', 1)),
	'POUND_KEY': createTokenizer(createMatch('#', 1)),
	'SEMICOLON': createTokenizer(createMatch(';', 1)),
	'AT_SYMBOL': createTokenizer(createMatch('@', 1)),
	'COMMA': createTokenizer(createMatch(',', 1)),
	'EQUAL_SYMBOL': createTokenizer(createMatch('=', 1)),
	'UNDERLINE': createTokenizer(createMatch('_', 1)),
	'PERSCENT_SYMBOL': createTokenizer(createMatch('%', 1)),
	'ASTERISK': createTokenizer(createMatch('*', 1)),
	'LEFT_ANGLE_BRACKET': createTokenizer(createMatch('<', 1)),
	'RIGHT_ANGLE_BRACKET': createTokenizer(createMatch('>', 1)),
	'UNKNOW': createTokenizer({
		start: () => true,
		end: () => true,
	})
}

function getToken(input, start) {
	for (let key in patterns) {
		let tokenizer = patterns[key]
		let token = tokenizer(input, start)
		if (token) {
			return Object.assign({type: key}, token)
		}
	}
}

function tokenizer(input) {
	let index = 0
	let tokens = []
	while (index < input.length) {
		let token = getToken(input, index)
		if (token) {
			tokens[tokens.length] = token
			index = token.end
		}
	}
	return tokens
}


const patterns1 = {
	'RULE': createTokenizer({
		last: true,
		getInitialValue: () => [],
		check: value => !!value,
		start: token => token.type === 'NAME',
		end: (token, value) => {
			return token.type === 'SEMICOLON'
		},
		getValue: (token, value, offset) => {
			if (value.length === 0) {
				value.push(token)
				return value
			}

			let last = value[value.length - 1]
			if (value.length === 1 && last.type === 'NAME' && token.type !== 'COLON') {
				return
			}

			value.push(token)
			return value
		},
		getFinalValue: tokens => tokens.map(token => token.value).join(''),
	}),
	// 'RULES': createTokenizer({
	// 	test: true,
	// 	last: true,
	// 	start: token => token.type === 'LEFT_BRACE',
	// 	end: token => token.type === 'RIGHT_BRACE',
	// 	getValue: (item, value) => value + item.value,
	// }),
	'SELECTOR': createTokenizer({
		start: token => token.type !== 'LEFT_BRACE' && token.type !== 'WHITE_SPACE' && token.type !== 'SEMICOLON',
		end: token => token.type === 'WHITE_SPACE' || token.type === 'LEFT_BRACE',
		getValue: (item, value) => value + item.value,
	}),
	'UNKNOW': createTokenizer({
		start: () => true,
		end: () => true,
		getValue: (item) => item,
	}),
}

function getToken1(input, start) {
	for (let key in patterns1) {
		let tokenizer = patterns1[key]
		let token = tokenizer(input, start)
		if (token) {
			return Object.assign({type: key}, token)
		}
	}
}

function tokenizer1(input) {
	let index = 0
	let tokens = []
	while (index < input.length) {
		let token = getToken1(input, index)
		if (token) {
			tokens[tokens.length] = token
			index = token.end
		}
	}
	return tokens
}




let fs = require('fs')
let path = require('path')
let cssFilePath = path.join(__dirname, 'files/test.css')
let content = fs.readFileSync(cssFilePath).toString()
let tokens = tokenizer(content)
let tokens1 = tokenizer1(tokens)

let destPath = path.join(__dirname, 'dest/05.json')
fs.writeFileSync(destPath, JSON.stringify(tokens1, null, 2))