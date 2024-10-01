/* eslint-disable no-console */
/* eslint-disable class-methods-use-this */

import { RuntimeError } from './consts'
import {
  Binary,
  Expr,
  Grouping,
  Literal,
  Unary,
  ExprVisitor,
  StmtVisitor,
  Expression,
  Print,
  Stmt
} from './expr'
import { IToken, TokenType } from './types'

export class Interpreter implements ExprVisitor<unknown>, StmtVisitor<void> {
  private readonly _reportError: (error: RuntimeError) => void

  constructor(options: { reportError?: (error: RuntimeError) => void } = {}) {
    this._reportError = options.reportError || console.error
  }

  interpret(statements: Stmt[]): void {
    try {
      for (const statement of statements) {
        this._execute(statement)
      }
    } catch (error) {
      if (error instanceof RuntimeError) {
        this._reportError(error)
      }
    }
  }

  visitExpressionStmt(expression: Expression): void {
    this._evaluate(expression.expression)
  }

  visitPrintStmt(print: Print): void {
    const value = this._evaluate(print.expression)
    console.log(this._stringify(value))
  }

  visitBinaryExpr(binary: Binary): any {
    const left = this._evaluate(binary.left)
    const right = this._evaluate(binary.right)

    switch (binary.operator.type) {
      case TokenType.BANG_EQUAL:
        return !this._isEqual(left, right)
      case TokenType.EQUAL_EQUAL:
        return this._isEqual(left, right)
      case TokenType.GREATER:
        this._checkNumberOperand2(binary.operator, left, right)
        return left > right
      case TokenType.GREATER_EQUAL:
        this._checkNumberOperand2(binary.operator, left, right)
        return left >= right
      case TokenType.LESS:
        this._checkNumberOperand2(binary.operator, left, right)
        return left < right
      case TokenType.LESS_EQUAL:
        this._checkNumberOperand2(binary.operator, left, right)
        return left <= right
      case TokenType.MINUS:
        this._checkNumberOperand2(binary.operator, left, right)
        return left - right
      case TokenType.PLUS:
        if (typeof left === 'number' && typeof right === 'number') return left + right
        if (typeof left === 'string' && typeof right === 'string') return left + right
        throw new RuntimeError(binary.operator, 'Operands must be two numbers or two strings.')
      case TokenType.SLASH:
        this._checkNumberOperand2(binary.operator, left, right)
        return left / right
      case TokenType.STAR:
        this._checkNumberOperand2(binary.operator, left, right)
        return left * right
      default:
        return undefined
    }
  }

  visitGroupingExpr(grouping: Grouping): any {
    return this._evaluate(grouping.expression)
  }

  visitLiteralExpr(literal: Literal): any {
    return literal.value
  }

  visitUnaryExpr(unary: Unary): any {
    const right = this._evaluate(unary.right)
    switch (unary.operator.type) {
      case TokenType.BANG:
        return !this._isTruthy(right)
      case TokenType.MINUS:
        this._checkNumberOperand1(unary.operator, right)
        return -right
      default:
        return undefined
    }
  }

  private _execute(stmt: Stmt): void {
    stmt.accept(this)
  }

  private _evaluate(expr: Expr): any {
    return expr.accept(this)
  }

  /**
   * @param operator 运算符.
   * @param operand 运算数.
   */
  private _checkNumberOperand1(operator: IToken, operand: unknown): void {
    if (typeof operand === 'number') return
    throw new RuntimeError(operator, 'Operand must be a number.')
  }

  private _checkNumberOperand2(operator: IToken, left: unknown, right: unknown): void {
    if (typeof left === 'number' && typeof right === 'number') return
    throw new RuntimeError(operator, 'Operands must be numbers.')
  }

  private _isTruthy(value: unknown): boolean {
    if (value == null) return false
    if (typeof value === 'boolean') return value
    return true
  }

  private _isEqual(a: unknown, b: unknown): boolean {
    return a === b
  }

  private _stringify(value: unknown): string {
    if (value == null) return 'nil'
    if (typeof value === 'number') {
      const res = value.toString()
      if (res.endsWith('.0')) return res.slice(0, -2)
      return res
    }
    return value.toString()
  }
}
