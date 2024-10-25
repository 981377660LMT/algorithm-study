import type { Node, ConditionalOperationNode, FunctionCallNode } from './ast'
import type { Parser } from './parser'

interface PrefixParseLet {
  handle(token: string, parser: Parser): Node
}

interface InfixParseLet {
  handle(left: Node, token: string, parser: Parser): Node
  precedence: number
}

export const prefixParselets: Record<string, PrefixParseLet> = {
  __value: {
    handle(token, parser) {
      return {
        type: 'value',
        value: token
      }
    }
  }
}

export const infixParselets: Record<string, InfixParseLet> = {}

// !JS运算符优先级(18种)
// use operator precedence of JavaScript
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Operator_Precedence

helpCreatePrefixOperator('+', 150)
helpCreatePrefixOperator('-', 150)
registerParenthesis()

helpCreateInfixOperator('+', 120)
helpCreateInfixOperator('-', 120)
helpCreateInfixOperator('*', 130)
helpCreateInfixOperator('/', 130)
helpCreateInfixOperator('^', 140, true)
registerConditionalOperater()
registerFunctionCall()

function helpCreatePrefixOperator(prefix: string, precedence: number) {
  prefixParselets[prefix] = {
    handle(token, { parseExp }) {
      const body = parseExp(precedence)
      return {
        type: 'unary',
        operator: prefix,
        right: body
      }
    }
  }
}

function helpCreateInfixOperator(infix: string, precedence: number, associateRight2Left = false) {
  infixParselets[infix] = {
    precedence,
    handle(left, token, { parseExp }) {
      const right = parseExp(associateRight2Left ? precedence - 1 : precedence)
      return {
        type: 'binary',
        operator: infix,
        left,
        right
      }
    }
  }
}

/**
 * Conditional operater is a special "infix"
 */
function registerConditionalOperater() {
  infixParselets['?'] = {
    // the binding power between condition node and "?"
    precedence: 30,
    handle(left, token, { parseExp, scanner }): ConditionalOperationNode {
      // the binding power to true/false branch is smallest
      // (parse expression as long as possible)
      const trueBranch = parseExp(0)
      scanner.consume(':')
      const falseBranch = parseExp(0)

      return {
        type: 'conditional',
        condition: left,
        trueBranch,
        falseBranch
      }
    }
  }
}

function registerParenthesis() {
  prefixParselets['('] = {
    handle(token, { parseExp, scanner }) {
      const content = parseExp(0)
      scanner.consume(')')
      return content
    }
  }
}

function registerFunctionCall() {
  infixParselets['('] = {
    precedence: 180,
    handle(left, token, { parseExp, scanner }): FunctionCallNode {
      const args: Node[] = []
      while (scanner.peek() !== ')') {
        // has arguments
        args.push(parseExp(0))
        if (scanner.peek() !== ',') break
        scanner.consume(',')
      }
      scanner.consume(')')
      return {
        type: 'fnCall',
        function: left,
        args
      }
    }
  }
}
