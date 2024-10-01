/* eslint-disable no-console */
/* eslint-disable class-methods-use-this */

import { createToken } from '../token'
import { TokenType } from '../types'
import { Visitor, Expr, Binary, Grouping, Literal, Unary } from './Expr'

/** @deprecated */
class AstPrinter implements Visitor<string> {
  visitBinaryExpr(binary: Binary): string {
    return this._parenthesize(binary.operator.lexeme, binary.left, binary.right)
  }

  visitGroupingExpr(grouping: Grouping): string {
    return this._parenthesize('group', grouping.expression)
  }

  visitLiteralExpr(literal: Literal): string {
    if (literal.value === undefined) return 'nil'
    return literal.value.toString()
  }

  visitUnaryExpr(unary: Unary): string {
    return this._parenthesize(unary.operator.lexeme, unary.right)
  }

  private _parenthesize(name: string, ...exprs: Expr[]): string {
    const sb: string[] = []
    sb.push('(')
    sb.push(name)
    exprs.forEach(v => {
      sb.push(' ')
      sb.push(v.accept(this))
    })
    sb.push(')')
    return sb.join('')
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const expr = new Binary(
    new Unary(
      createToken({ type: TokenType.MINUS, lexeme: '-', literal: null, line: 1 }),
      new Literal(123)
    ),
    createToken({ type: TokenType.STAR, lexeme: '*', literal: null, line: 1 }),
    new Grouping(new Literal(45.67))
  )

  const astPrinter = new AstPrinter()
  console.log(astPrinter.visitBinaryExpr(expr))
}
