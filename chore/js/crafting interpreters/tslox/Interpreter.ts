/* eslint-disable no-console */
/* eslint-disable class-methods-use-this */

import { Return, RuntimeError } from './consts'
import { Environment } from './Environment'
import { LoxCallable, LoxFunction, LoxClass, LoxInstance } from './callable'
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
  Stmt,
  VariableExpr,
  VariableDecl,
  Assign,
  Block,
  IfStmt,
  Logical,
  WhileStmt,
  Call,
  Func,
  ReturnStmt,
  ClassStmt,
  Get,
  SetExpr,
  ThisExpr,
  SuperExpr
} from './syntaxTree'
import { IToken, TokenType } from './types'

export class Interpreter implements ExprVisitor<unknown>, StmtVisitor<void> {
  readonly globalEnv = new Environment()
  private _env = this.globalEnv

  /** The distance of a variable from the current scope. */
  private readonly _locals = new WeakMap<Expr, number>()

  private readonly _reportError: (error: RuntimeError) => void

  constructor(options: { reportError?: (error: RuntimeError) => void } = {}) {
    this._reportError = options.reportError || console.error

    /** Native function clock(). */
    this.globalEnv.define(
      'clock',
      new (class extends LoxCallable {
        override arity(): number {
          return 0
        }

        override call(): number {
          return Date.now()
        }
      })()
    )
  }

  interpret(statements: Stmt[]): void {
    try {
      for (const statement of statements) {
        this._execute(statement)
      }
    } catch (error) {
      if (error instanceof RuntimeError) {
        this._reportError(error)
      } else {
        throw error
      }
    }
  }

  visitIfStmtStmt(ifstmt: IfStmt): void {
    if (this._isTruthy(this._evaluate(ifstmt.condition))) {
      this._execute(ifstmt.thenBranch)
    } else if (ifstmt.elseBranch) {
      this._execute(ifstmt.elseBranch)
    }
  }

  visitExpressionStmt(expression: Expression): void {
    this._evaluate(expression.expression)
  }

  visitFuncStmt(func: Func): void {
    const f = new LoxFunction().init(func, this._env)
    this._env.define(func.name.lexeme, f)
  }

  visitReturnStmtStmt(returnStmt: ReturnStmt): void {
    let value: any
    if (returnStmt.value != undefined) {
      value = this._evaluate(returnStmt.value)
    }
    throw new Return(value)
  }

  visitPrintStmt(print: Print): void {
    const value = this._evaluate(print.expression)
    console.log(this._stringify(value))
  }

  visitBlockStmt(block: Block): void {
    this.executeBlock(block.statements, new Environment(this._env))
  }

  visitClassStmtStmt(classstmt: ClassStmt): void {
    let superclass: LoxClass | undefined
    if (classstmt.superclass) {
      superclass = this._evaluate(classstmt.superclass)
      if (!(superclass instanceof LoxClass)) {
        throw new RuntimeError(classstmt.superclass.name, 'Superclass must be a class.')
      }
    }

    this._env.define(classstmt.name.lexeme, undefined)

    if (classstmt.superclass) {
      this._env = new Environment(this._env)
      this._env.define('super', superclass)
    }

    const methods = new Map<string, LoxFunction>()
    classstmt.methods.forEach(method => {
      const f = new LoxFunction().init(method, this._env, method.name.lexeme === 'init')
      methods.set(method.name.lexeme, f)
    })
    const cls = new LoxClass(classstmt.name.lexeme, superclass, methods)

    if (superclass) {
      this._env = this._env.enclosing!
    }
    this._env.assign(classstmt.name, cls)
  }

  visitVariableDeclStmt(variableDecl: VariableDecl): void {
    let value: any
    if (variableDecl.initializer != undefined) {
      value = this._evaluate(variableDecl.initializer)
    }
    this._env.define(variableDecl.name.lexeme, value)
  }

  visitWhileStmtStmt(whilestmt: WhileStmt): void {
    while (this._isTruthy(this._evaluate(whilestmt.condition))) {
      this._execute(whilestmt.body)
    }
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

  visitCallExpr(call: Call): unknown {
    const callee = this._evaluate(call.callee)
    if (!(callee instanceof LoxCallable)) {
      throw new RuntimeError(call.paren, 'Can only call functions and classes.')
    }
    const args = call.args.map(arg => this._evaluate(arg))
    if (args.length !== callee.arity()) {
      throw new RuntimeError(
        call.paren,
        `Expected ${callee.arity()} arguments but got ${args.length}.`
      )
    }
    return callee.call(this, args)
  }

  visitGetExpr(get: Get): unknown {
    const obj = this._evaluate(get.obj)
    if (obj instanceof LoxInstance) {
      return obj.get(get.name)
    }
    throw new RuntimeError(get.name, 'Only instances have properties.')
  }

  visitGroupingExpr(grouping: Grouping): any {
    return this._evaluate(grouping.expression)
  }

  visitLiteralExpr(literal: Literal): any {
    return literal.value
  }

  visitLogicalExpr(logical: Logical): unknown {
    const left = this._evaluate(logical.left)
    if (logical.operator.type === TokenType.OR) {
      if (this._isTruthy(left)) return left
    } else if (!this._isTruthy(left)) return left

    // We look at its value to see if we can short-circuit.
    // If not, and only then, do we evaluate the right operand.
    return this._evaluate(logical.right)
  }

  visitSetExprExpr(setexpr: SetExpr): unknown {
    const obj = this._evaluate(setexpr.obj)
    if (!(obj instanceof LoxInstance)) {
      throw new RuntimeError(setexpr.name, 'Only instances have fields.')
    }

    const value = this._evaluate(setexpr.value)
    obj.set(setexpr.name, value)
    return value
  }

  visitSuperExprExpr(superexpr: SuperExpr): unknown {
    const dist = this._locals.get(superexpr)!
    const superclass = this._env.getAt(dist, 'super') as LoxClass
    const obj = this._env.getAt(dist - 1, 'this') as LoxInstance
    const method = superclass.findMethod(superexpr.method.lexeme)
    if (!method) {
      throw new RuntimeError(superexpr.method, `Undefined property '${superexpr.method.lexeme}'.`)
    }
    return method.bind(obj)
  }

  visitThisExprExpr(thisexpr: ThisExpr): unknown {
    return this._lookUpVariable(thisexpr.keyword, thisexpr)
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

  visitVariableExprExpr(variableExpr: VariableExpr): unknown {
    return this._lookUpVariable(variableExpr.name, variableExpr)
  }

  visitAssignExpr(assign: Assign): unknown {
    const value = this._evaluate(assign.value)
    const dist = this._locals.get(assign)
    if (dist != undefined) {
      this._env.assignAt(dist, assign.name, value)
    } else {
      this.globalEnv.assign(assign.name, value)
    }
    return value
  }

  private _execute(stmt: Stmt): void {
    stmt.accept(this)
  }

  resolve(expr: Expr, depth: number): void {
    this._locals.set(expr, depth)
  }

  /** Executes a list of statements in the context of a given environment. */
  executeBlock(statements: Stmt[], env: Environment): void {
    const preEnv = this._env
    try {
      this._env = env
      for (const statement of statements) {
        this._execute(statement)
      }
    } finally {
      this._env = preEnv
    }
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
    if (value == undefined) return false
    if (typeof value === 'boolean') return value
    return true
  }

  private _isEqual(a: unknown, b: unknown): boolean {
    return a === b
  }

  private _stringify(value: unknown): string {
    if (value == undefined) return 'nil'
    if (typeof value === 'number') {
      const res = value.toString()
      if (res.endsWith('.0')) return res.slice(0, -2)
      return res
    }
    return value.toString()
  }

  private _lookUpVariable(name: IToken, expr: Expr): any {
    const dist = this._locals.get(expr)
    if (dist != undefined) {
      return this._env.getAt(dist, name.lexeme)
    }
    return this.globalEnv.get(name)
  }
}
