const isWhiteSpace = character => {
	switch (character) {
		case ' ':
		case '\t':
		case '\r':
			return true
		default:
			return false
	}
}

const isOperand = character => {
	switch (character) {
		case '+':
		case '-':
		case '*':
		case '/':
			return true
		default:
			return false
	}
}

const isDigit = character => {
	return character >= '0' && character <= '9'
}

const isParenthese = character => {
	return character === '(' || character === ')'
}

const tokenizer = input => {
	let tokenList = []
	let index = 0
	let character

	let next = () => {
		character = input[index++]
	}

	let peek = (n = 0) => {
		return input[index + n]
	}

	let ignoreWhiteSpace = () => {
		while (isWhiteSpace(character)) {
			next()
		}
	}

	let consume = () => {
		next()
		ignoreWhiteSpace()
	}

	let handleNumber = () => {
		let value = character
		let nextCharacter = peek()
		while (isDigit(nextCharacter) || nextCharacter === '.') {
			value += nextCharacter
			next()
			nextCharacter = peek()
		}
		tokenList.push({
			type: 'number',
			value: value
		})
	}

	let handleOperand = () => {
		let isNumber = (character === '-' || character === '+') && isDigit(peek())
		if (isNumber) {
			return handleNumber()
		}
		tokenList.push({
			type: 'operand',
			value: character
		})
	}

	let handleParenthese = () => {
		tokenList.push({
			type: 'parenthese',
			value: character
		})
	}

	while (index < input.length) {
		consume()

		if (isOperand(character)) {
			handleOperand()
			continue
		}

		if (isDigit(character)) {
			handleNumber()
			continue
		}

		if (isParenthese(character)) {
			handleParenthese()
			continue
		}
	}

	return tokenList
}

const Program = 'Program'
const ExpressionStatement = 'ExpressionStatement'
const NumberLiteral = 'NumberLiteral'
const GroupingStatement = 'GroupingStatement'

function parser(input) {
	let ast = {
		type: Program,
		body: []
	}
	let stack = ast.body
	let index = 0
	let token

	let next = () => {
		token = input[index++]
	}

	let peek = () => {
		return input[index]
	}

	let handleRight = (previous, right) => {
		let target = previous
		while (target.right && target.right.type === ExpressionStatement) {
			target = target.right
		}
		target.right = right
	}

	let handleNumber = () => {
		let previous = stack[stack.length - 1]
		if (!previous) {
			stack.push({
				type: NumberLiteral,
				value: token.value
			})
			return
		}

		if (previous.type !== ExpressionStatement) {
			throw new Error('非法组合：数字前面必须是操作符')
		}

		handleRight(previous, {
			type: NumberLiteral,
			value: token.value
		})
	}

	let handleExpression = () => {
		let previous = stack[stack.length - 1]
		if (!previous) {
			throw new Error('操作符不能单独出现')
		}

		let isValid =
			previous.type === NumberLiteral ||
			previous.type === ExpressionStatement ||
			previous.type === GroupingStatement

		if (!isValid) {
			throw new Error('非法表达式')
		}

		if (token.value === '+' || token.value === '-') {
			stack.pop()
			let expression = {
				type: ExpressionStatement,
				operand: token.value,
				left: previous,
				right: null
			}
			stack.push(expression)
			return
		} else if (token.value === '*' || token.value === '/') {
			if (previous.type === ExpressionStatement) {
				if (previous.operand === '+' || previous.operand === '-') {
					let expression = {
						type: ExpressionStatement,
						operand: token.value,
						left: previous.right,
						right: null
					}
					previous.right = expression
					return
				} else if (previous.operand === '*' || previous.operand === '/') {
					stack.pop()
					let expression = {
						type: ExpressionStatement,
						operand: token.value,
						left: previous,
						right: null
					}
					stack.push(expression)
					return
				}
			} else if (
				previous.type === NumberLiteral ||
				previous.type === GroupingStatement
			) {
				stack.pop()
				let expression = {
					type: ExpressionStatement,
					operand: token.value,
					left: previous,
					right: null
				}
				stack.push(expression)
				return
			}
		}

		throw new Error('不支持的操作')
	}

	let handleGrouping = () => {
		let grouping = {
			type: GroupingStatement,
			body: []
		}

		let currentStack = stack
		stack = grouping.body

		let currentToken = token
		while (
			currentToken &&
			!(currentToken.type === 'parenthese' && currentToken.value === ')')
		) {
			handler()
			currentToken = peek()
		}

		next()

		stack = currentStack

		if (!grouping.body.length) {
			throw new Error('不支持空括号')
		}

		let previous = stack[stack.length - 1]
		if (!previous) {
			stack.push(grouping)
		}

		if (previous.type === ExpressionStatement) {
			handleRight(previous, grouping)
			return
		}

		throw new Error('分组前面的符号非法')
	}

	let handler = () => {
		next()
		switch (token && token.type) {
			case 'number':
				handleNumber()
				break
			case 'operand':
				handleExpression()
				break
			case 'parenthese':
				handleGrouping()
				break
		}
	}

	while (index < input.length) {
		handler()
	}

	return ast
}

function traverser(ast, visitor) {
	let handleArray = (array, parent, layer) => {
		array.forEach(node => handleNode(node, parent, layer))
	}
	let handleNode = (node, parent, layer) => {
		let enter = `on${node.type}Enter`
		let through = `on${node.type}Through`
		let exit = `on${node.type}Exit`

		if (visitor[enter]) {
			visitor[enter](node, parent, layer)
		}

		switch (node.type) {
			case Program:
				handleArray(node.body, node, layer + 1)
				break
			case ExpressionStatement:
				handleNode(node.left, node, layer + 1)
				if (visitor[through]) {
					visitor[through](node, parent, layer)
				}
				handleNode(node.right, node, layer + 1)
				break
			case GroupingStatement:
				handleArray(node.body, node, layer + 1)
				break
			case NumberLiteral:
				break
			default:
				throw new TypeError(node.type)
		}

		if (visitor[exit]) {
			visitor[exit](node, parent, layer)
		}
	}
	handleNode(ast, null, 0)
}

function toXML(ast, indent = 2) {
	let xml = ''
	let handleIndent = (layer = 0) => {
		if (indent > 0) {
			return '\n'.padEnd(layer * indent, ' ')
		}
		return ''
	}
	let visitor = {
		onProgramEnter(node, parent, layer) {
			xml += handleIndent(layer) + '<program>'
		},
		onProgramExit(node, parent, layer) {
			xml += handleIndent(layer) + '</program>'
		},
		onNumberLiteralEnter(node, parent, layer) {
			xml += handleIndent(layer) + '<number>'
			xml += handleIndent(layer + 1) + node.value
		},
		onNumberLiteralExit(node, parent, layer) {
			xml += handleIndent(layer) + '</number>'
		},
		onExpressionStatementEnter(node, parent, layer) {
			xml += handleIndent(layer) + `<expression operand="${node.operand}">`
		},
		onExpressionStatementExit(node, parent, layer) {
			xml += handleIndent(layer) + '</expression>'
		},
		onGroupingStatementEnter(node, parent, layer) {
			xml += handleIndent(layer) + '<grouping>'
		},
		onGroupingStatementExit(node, parent, layer) {
			xml += handleIndent(layer) + '</grouping>'
		}
	}
	traverser(ast, visitor)
	return xml
}

function toChinese(ast) {
	let string = ''
	let visitor = {
		onExpressionStatementThrough(node) {
			switch (node.operand) {
				case '+':
					string += ' 加 '
					break
				case '-':
					string += ' 减 '
					break
				case '*':
					string += ' 乘 '
					break
				case '/':
					string += ' 除 '
					break
				default:
					throw new Error(`unknow operand ${node.operand}`)
			}
		},
		onNumberLiteralExit(node) {
			let value = node.value
				.split('')
				.map(item => item === '.' ? '点' : '零一二三四五六七八九'[item])
				.join('')
			string += value
		},
		onGroupingStatementEnter() {
			string += '('
		},
		onGroupingStatementExit() {
			string += ')'
		}
	}
	traverser(ast, visitor)
	return string
}

function evaluate(ast) {
	let stack = []
	let visitor = {
		onExpressionStatementExit(node) {
			let right = stack.pop()
			let left = stack.pop()
			let value
			switch (node.operand) {
				case '+':
					value = left + right
					break
				case '-':
					value = left - right
					break
				case '*':
					value = left * right
					break
				case '/':
					value = left / right
					break
				default:
					throw new Error(`unknow operand ${node.operand}`)
			}
			stack.push(value)
		},
		onNumberLiteralEnter(node) {
			stack.push(Number(node.value))
		}
	}
	traverser(ast, visitor)
	return stack[0]
}

let expression = `1 + 2 + 3 * 4 + (5 + 6 * (7 + 8)) - 9/ 100 * 2 + 0.5 + -1.5`
let tokenList = tokenizer(expression)
console.log('tokenList', tokenList)

let ast = parser(tokenList)
console.log('ast', ast)

let xml = toXML(ast)
console.log('xml', xml)

let result = evaluate(ast)
console.log('result', result)

let chinese = toChinese(ast)
console.log('chinese', chinese)