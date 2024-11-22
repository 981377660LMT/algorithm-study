const SPACE = ' '
const TAB = '\t'
const NEW_LINE = '\n'
const ZERO = '0'
const NINE = '9'
const DOT = '.'
const PLUS_SIGN = '+'
const MINUS_SIGN = '-'
const MULTIPLICATION_SIGN = '*'
const DEVISION_SIGN = '/'
const LEFT_PARENTHESE = '('
const RIGHT_PARENTHESE = ')'

const NUMBER = 0
const OPERAND = 1
const UNKNOW = 2
const PARENTHESE = 3
const GROUPING = 4

class Character {
	constructor(value) {
		this.value = value
	}
	isWhiteSpace() {
		switch (this.value) {
			case SPACE:
			case TAB:
			case NEW_LINE:
				return true
			default:
				return false
		}
	}
	isOperand() {
		switch (this.value) {
			case PLUS_SIGN:
			case MINUS_SIGN:
			case MULTIPLICATION_SIGN:
			case DEVISION_SIGN:
				return true
			default:
				return false
		}
	}
	isDigit() {
		return this.value >= ZERO && this.value <= NINE
	}
	isDot() {
		return this.value === DOT
	}
	isLeftParenthese() {
		return this.value === LEFT_PARENTHESE
	}
	isRightParenthese() {
		return this.value === RIGHT_PARENTHESE
	}
}

class Token {
	constructor(type, startIndex) {
		this.type = type
		this.startIndex = startIndex
		this.value = ''
	}
	add(character) {
		this.value += character
	}
	is(type) {
		return this.type === type
	}
}

class Tokenizer {
	constructor(input, output = []) {
		this.input = input
		this.output = output
		this.index = 0
		this.character = null
	}
	isNotEnd() {
		return this.index < this.input.length
	}
	next() {
		if (this.isNotEnd()) {
			return new Character(this.input[this.index++])
		} else {
			return null
		}
	}
	peek(n = 0) {
		if (this.index + n < this.input.length) {
			return new Character(this.input[this.index + n])
		} else {
			return null
		}
	}
	consume() {
		let character = this.next()
		while (character && character.isWhiteSpace()) {
			character = this.next()
		}
		this.character = character
		return !!this.character
	}
	createToken(type) {
		return new Token(type, this.index - 1)
	}
	handleNumber() {
		let token = this.createToken(NUMBER)
		token.add(this.character.value)

		let nextToken = this.peek()

		while (nextToken && (nextToken.isDigit() || nextToken.isDot())) {
			token.add(this.next().value)
			nextToken = this.peek()
		}

		if (token.value.split('.').filter(item => item === DOT).length > 1) {
			throw new Error('非法数字')
		}

		this.output.push(token)
	}
	handleOperand() {
		let token = this.createToken(OPERAND)
		token.add(this.character.value)
		this.output.push(token)
	}
	handleParenthese() {
		let token = this.createToken(PARENTHESE)
		token.add(this.character.value)
		this.output.push(token)
	}
	execute() {
		while (this.isNotEnd()) {
			if (!this.consume()) {
				return this.output
			}

			let character = this.character

			let isDigit = character.isDigit()
			if (isDigit) {
				this.handleNumber()
				continue
			}

			let isOperand = character.isOperand()
			let isDigitWithOperand =
				(character.value === PLUS_SIGN || character.value === MINUS_SIGN) &&
				this.peek() &&
				this.peek().isDigit()

			if (isDigitWithOperand) {
				this.handleNumber()
				continue
			}

			if (isOperand) {
				this.handleOperand()
				continue
			}

			let isParenthese =
				character.isLeftParenthese() || character.isRightParenthese()
			if (isParenthese) {
				this.handleParenthese()
				continue
			}
		}

		return this.output
	}
}

class Ast {
	constructor() {
		this.type = 'Program'
		this.body = []
	}
	addNode(node) {
		this.body.push(node)
	}
}

class ExpressionStatement {
	constructor() {
		this.type = 'ExpressionStatement'
		this.operand = null
		this.left = null
		this.right = null
	}
	setOperand(operand) {
		this.operand = operand
	}
	setLeft(left) {
		this.left = left
	}
	setRight(right) {
		this.right = right
	}
}

class GroupingStatement {
	constructor() {
		this.type = 'GroupingStatement'
		this.body = []
	}
	addNode(node) {
		this.body.push(node)
	}
}

class NumberLiteral {
	constructor(value) {
		this.type = 'NumberLiteral'
		this.value = value
	}
}

class Parser {
	constructor(input, output = []) {
		this.input = input
		this.output = new Ast()
		this.index = 0
	}
	isNotEnd() {
		return this.index < this.input.length
	}
	next() {
		if (this.isNotEnd()) {
			return this.input[this.index++]
		} else {
			null
		}
	}
	peek(n = 0) {
		if (this.index + n < this.input.length) {
			return this.input[this.index + n]
		} else {
			return null
		}
	}
	match(...args) {
		return args.every(
			(type, index) => (this.peek(index) ? this.peek(index).is(type) : false)
		)
	}
	handleNumberLiteral() {
		let token = this.next()
		let previous = this.stack[this.stack.length - 1]

		if (!previous) {
			this.stack.push(new NumberLiteral(token.value))
			return
		}

		if (!(previous instanceof ExpressionStatement)) {
			throw new Error('非法组合：数字前面必须是操作符')
		}

		if (previous instanceof ExpressionStatement) {
			let target = previous
			while (target.right instanceof ExpressionStatement) {
				target = target.right
			}
			target.right = new NumberLiteral(token.value)
			return
		}
	}

	handleExpressionStatement() {
		let token = this.next()
		let previous = this.stack[this.stack.length - 1]

		if (!previous) {
			throw new Error('操作符不能单独出现')
		}

		let isValid =
			previous instanceof NumberLiteral ||
			previous instanceof ExpressionStatement ||
			previous instanceof GroupingStatement

		if (!isValid) {
			throw new Error(
				`操作符前面的类型必须是数字，表达式或者分组，而不是${previous}`
			)
		}

		if (token.value === PLUS_SIGN || token.value === MINUS_SIGN) {
			this.stack.pop()
			let expression = new ExpressionStatement()
			expression.setOperand(token.value)
			expression.setLeft(previous)
			this.stack.push(expression)
			return
		}

		if (token.value === MULTIPLICATION_SIGN || token.value === DEVISION_SIGN) {
			if (previous instanceof ExpressionStatement) {
				// example: 1 + (2 + 1) * 3 + 1 * 2 * 3
				if (previous.operand === PLUS_SIGN || previous.operand === MINUS_SIGN) {
					let expression = new ExpressionStatement()
					expression.setOperand(token.value)
					expression.setLeft(previous.right)
					previous.setRight(expression)
					return
				}

				// example: 1 * 2 * 3
				this.stack.pop()
				let expression = new ExpressionStatement()
				expression.setOperand(token.value)
				expression.setLeft(previous)
				this.stack.push(expression)
				return
			}

			if (
				previous instanceof NumberLiteral ||
				previous instanceof GroupingStatement
			) {
				this.stack.pop()
				let expression = new ExpressionStatement()
				expression.setOperand(token.value)
				expression.setLeft(previous)
				this.stack.push(expression)
				return
			}

			throw new Error('不支持的操作')
		}
	}

	handleGroupingStatement() {
		let token = this.next()
		let previous = this.stack[this.stack.length - 1]
		let currentStack = this.stack
		let grouping = new GroupingStatement()
		this.stack = grouping.body
		while (
			!(this.match(PARENTHESE) && this.peek().value === RIGHT_PARENTHESE)
		) {
			this.handle()
		}
		this.next()
		this.stack = currentStack
		if (!grouping.body.length) {
			throw new Error('不支持空括号')
		}

		if (!previous) {
			this.stack.push(grouping)
			return
		}

		if (previous instanceof ExpressionStatement) {
			let target = previous
			while (target.right instanceof ExpressionStatement) {
				target = target.right
			}
			target.setRight(grouping)
			return
		}

		throw new Error('分组前面的符号非法')
	}

	// 1 * 2 + 2 + (3 + 4) * 2 * 3 / 2
	handle() {
		if (this.match(NUMBER)) {
			return this.handleNumberLiteral()
		}

		if (this.match(OPERAND)) {
			return this.handleExpressionStatement()
		}

		if (this.match(PARENTHESE)) {
			return this.handleGroupingStatement()
		}

		throw new Error(`非法 token：${this.next()}`)
	}
	execute() {
		while (this.isNotEnd()) {
			this.stack = this.output.body
			this.handle()
		}
		return this.output
	}
}

class Traverser {
	constructor(ast, visitor) {
		this.ast = ast
		this.visitor = visitor
	}
	handleArray(array, parent, layer) {
		array.forEach(node => this.handleNode(node, parent, layer))
	}
	handleNode(node, parent, layer) {
		let enterMethod = `handle${node.type}Enter`
		let exitMethod = `handle${node.type}Exit`

		if (this.visitor[enterMethod]) {
			this.visitor[enterMethod](node, parent, layer)
		}

		if (node instanceof Ast) {
			this.handleArray(node.body, node, layer + 1)
		}

		if (node instanceof ExpressionStatement) {
			this.handleNode(node.left, node, layer + 1)
			this.handleNode(node.right, node, layer + 1)
		}

		if (node instanceof GroupingStatement) {
			this.handleArray(node.body, node, layer + 1)
		}

		if (this.visitor[exitMethod]) {
			this.visitor[exitMethod](node, parent, layer)
		}
	}
	execute() {
		this.handleNode(this.ast, null, 0)
	}
}

class Visitor {
	constructor(ast) {
		this.traverser = new Traverser(ast, this)
		this.output = null
	}
	execute() {
		this.traverser.execute()
		return this.output
	}
}

class XMLPrinter extends Visitor {
	constructor(ast) {
		super(ast)
		this.output = ''
		this.layer = 0
	}
	padStart(layer) {
		return ''.padStart(layer * 2, ' ')
	}
	handleProgramEnter(node, parent, layer) {
		this.output += this.padStart(layer)
		this.output += `<program>\n`
	}
	handleProgramExit(node, parent, layer) {
		this.output += this.padStart(layer)
		this.output += `</program>\n`
	}
	handleExpressionStatementEnter(node, parent, layer) {
		this.output += this.padStart(layer)
		this.output += `<expression operand="${node.operand}">\n`
	}
	handleExpressionStatementExit(node, parent, layer) {
		this.output += this.padStart(layer)
		this.output += `</expression>\n`
	}
	handleNumberLiteralEnter(node, parent, layer) {
		this.output += this.padStart(layer)
		this.output += `<number>${node.value}</number>\n`
	}
	handleNumberLiteralExit(node, parent, layer) {}
	handleGroupingStatementEnter(node, parent, layer) {
		this.output += this.padStart(layer)
		this.output += `<grouping>\n`
	}
	handleGroupingStatementExit(node, parent, layer) {
		this.output += this.padStart(layer)
		this.output += `</grouping>\n`
	}
}

class Interpreter extends Visitor {
	constructor(ast) {
		super(ast)
		this.stack = []
	}
	get output() {
		return this.stack[0]
	}
	set output(v) {}
	handleProgramEnter(node, parent, layer) {}
	handleProgramExit(node, parent, layer) {}
	handleExpressionStatementEnter(node, parent, layer) {}
	handleExpressionStatementExit(node, parent, layer) {
		let right = this.stack.pop()
		let left = this.stack.pop()
		let value = 0
		switch (node.operand) {
			case PLUS_SIGN:
				value = left + right
				break
			case MINUS_SIGN:
				value = left - right
				break
			case MULTIPLICATION_SIGN:
				value = left * right
				break
			case DEVISION_SIGN:
				value = left / right
				break
		}
		this.stack.push(value)
	}
	handleNumberLiteralEnter(node, parent, layer) {
		this.stack.push(Number(node.value))
	}
	handleNumberLiteralExit(node, parent, layer) {}
	handleGroupingStatementEnter(node, parent, layer) {}
	handleGroupingStatementExit(node, parent, layer) {}
}

function test() {
	let expressiton = `1 + 2 + 3 * 4 + (5 + 6 * (7 + 8)) - 9/ 100 * 2 + 0.5 + -1.5`
	let tokenizer = new Tokenizer(expressiton)
	let parser = new Parser(tokenizer.execute())
	parser.execute()

	let xmlPrinter = new XMLPrinter(parser.output)
	let interpreter = new Interpreter(parser.output)

	xmlPrinter.execute()
	interpreter.execute()
	console.log('ast', parser.output)
	console.log('result', interpreter.output)
	console.log('xml')
	console.log(xmlPrinter.output)
	document.documentElement.innerHTML = xmlPrinter.output
}

test()
