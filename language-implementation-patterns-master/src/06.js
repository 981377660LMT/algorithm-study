const isWhiteSpace = char => (
	char === ' ' ||
	char === '\n' ||
	char === '\r' ||
	char === '\t'
)
const isLetter = char => char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'

module.exports = parse

function parse(inputs) {
	let groups = toGroup(inputs)
	let structure = toStructure(groups)
	return structure
}

function toGroup(inputs) {
	let index = 0
	let groups = []
	let char = ''
	let ignoreWhiteSpace = () => {
		while (isWhiteSpace(char)) {
			char = inputs[index++]
		}
	}
	let getNext = () => {
		char = inputs[index++]
	}

	getNext()

	while (char != null) {
		ignoreWhiteSpace()

		if (char == null) {
			break
		}

		let value = ''

		while (char === '`') {
			value += char
			getNext()
		}

		if (value) {
			groups.push({
				type: 'BACKQUOTE',
				value: value,
			})
			continue
		}

		while (isLetter(char)) {
			value += char
			getNext()
		}

		if (value) {
			groups.push({
				type: 'LETTER',
				value: value,
			})
			continue
		}

		if (char === ':') {
			groups.push({
				type: 'COLON',
				value: char,
			})
			getNext()
			continue
		}

		groups.push({
			type: 'UNKNOW',
			value: char,
		})
		getNext()
	}

	return groups
}

function toStructure(groups) {
	let structure = []
	let index = 0
	let currentItem = groups[0]

	let isBlockStart = () => {
		if (index >= groups.length) {
			return false
		}
		return (
			groups[index].type === 'BACKQUOTE' &&
			index + 1 < groups.length &&
			groups[index + 1].type === 'LETTER'
		)
	}

	let isPair = () => {
		if (index >= groups.length) {
			return false
		}
		return (
			groups[index].type === 'LETTER' &&
			index + 1 < groups.length &&
			groups[index + 1].type === 'COLON'
		)
	}

	let isBlockEnd = () => {
		if (index >= groups.length) {
			return false
		}
		return (
			groups[index].type === 'BACKQUOTE' &&
			groups[index].value.length === 3
		)
	}

	while (index < groups.length) {
		while (isBlockStart())  {
			let block = {
				type: groups[index + 1].value,
				data: {},
			}

			index += 2

			while (isPair()) {
				let key = groups[index].value
				let value = ''
				index += 2
				while (!isPair() && !isBlockEnd()) {
					if (index >= groups.length) {
						break
					}
					value += groups[index].value
					index += 1
				}
				block.data[key] = value
			}

			if (!isBlockEnd()) {
				throw new Error(`Expected \`\`\`, but get ${JSON.stringify(groups[index], null, 2)}`)
			}
			index += 1
			structure.push(block)
		}
	}
	return structure
}