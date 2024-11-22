var expect = require('expect')
var ListLexer = require('./01')

describe('test ListLexer', () => {
	it('should work as expected', () => {
		let lexer = new ListLexer('[a, b]')
		let tokens = []
		let token = lexer.nextToken()

		while (token.type !== ListLexer.EOF_TYPE) {
			tokens.push(token.toString())
			token = lexer.nextToken()
		}

		expect(tokens).toEqual([
			`<'[',LBRACK>`,
			`<'a',NAME>`,
			`<',',COMMA>`,
			`<'b',NAME>`,
			`<']',RBRACK>`
		])
	})
})