import type { Node } from './ast'
import { prefixParselets, infixParselets } from './parselet'
import type { Scanner } from './scanner'

export interface Parser {
  parseProgram(): Node
  parseExp(ctxPrecedence: number): Node
  scanner: Scanner
}

export function createParser(scanner: Scanner): Parser {
  const parser: Parser = {
    parseProgram,
    parseExp,
    scanner
  }
  return parser

  // for this naive parser,
  // a program is just one expression
  function parseProgram() {
    return parseExp(0)
  }

  function parseExp(ctxPrecedence: number): Node {
    let prefixToken = scanner.consume()
    if (!prefixToken) throw new Error('expect token but found none')

    // because our scanner is so naive,
    // we treat all non-operator tokens as value (.e.g number)
    const prefixParselet = prefixParselets[prefixToken] ?? prefixParselets.__value
    let left: Node = prefixParselet.handle(prefixToken, parser)

    while (true) {
      const infixToken = scanner.peek()
      if (!infixToken) break
      const infixParselet = infixParselets[infixToken]
      if (!infixParselet) break
      if (infixParselet.precedence <= ctxPrecedence) break
      scanner.consume()
      left = infixParselet.handle(left, infixToken, parser)
    }
    return left
  }
}
