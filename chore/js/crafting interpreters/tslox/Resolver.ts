/* eslint-disable no-console */
import { type Interpreter } from './Interpreter'
import {
  Assign,
  Binary,
  Block,
  Call,
  ClassStmt,
  Expr,
  Expression,
  ExprVisitor,
  Func,
  Get,
  Grouping,
  IfStmt,
  Literal,
  Logical,
  Print,
  ReturnStmt,
  SetExpr,
  Stmt,
  StmtVisitor,
  SuperExpr,
  ThisExpr,
  Unary,
  VariableDecl,
  VariableExpr,
  WhileStmt
} from './syntaxTree'
import { IToken, ReportErrorFunc } from './types'
import { noop } from './utils'

enum FunctionType {
  None,
  Function,
  Initializer,
  Method
}

enum ClassType {
  None,
  Class,
  SubClass
}

export class Resolver implements ExprVisitor<void>, StmtVisitor<void> {
  private readonly _interpreter: Interpreter
  private readonly _reportError: ReportErrorFunc

  /**
   * variable => whether it's been defined.
   */
  private readonly _scopes: Map<string, boolean>[] = []

  /** track whether or not the code we are currently visiting is inside a function declaration. */
  private _currentFunction = FunctionType.None
  private _currentClass = ClassType.None

  constructor(interpreter: Interpreter, options: { reportError?: ReportErrorFunc } = {}) {
    this._interpreter = interpreter
    this._reportError = options.reportError || console.error
  }

  visitBlockStmt(block: Block): void {
    this._beginScope()
    this.resolve(block.statements)
    this._endScope()
  }

  visitClassStmtStmt(classStmt: ClassStmt): void {
    const enclosingClass = this._currentClass
    this._currentClass = ClassType.Class
    this._declare(classStmt.name)
    this._define(classStmt.name)

    if (classStmt.superclass) {
      if (classStmt.name.lexeme === classStmt.superclass.name.lexeme) {
        this._reportError(classStmt.superclass.name, 'A class cannot inherit from itself.')
      }
      this._currentClass = ClassType.SubClass
      this._resolveExpr(classStmt.superclass)
      this._beginScope()
      this._scopes[this._scopes.length - 1].set('super', true)
    }

    this._beginScope()
    this._scopes[this._scopes.length - 1].set('this', true)
    classStmt.methods.forEach(method => {
      let declaration = FunctionType.Method
      if (method.name.lexeme === 'init') {
        declaration = FunctionType.Initializer
      }
      this._resolveFunction(method, declaration)
    })
    this._endScope()
    if (classStmt.superclass) {
      this._endScope()
    }
    this._currentClass = enclosingClass
  }

  /** split binding into two steps, declaring then defining. */
  visitVariableDeclStmt(variabledecl: VariableDecl): void {
    this._declare(variabledecl.name)
    if (variabledecl.initializer) {
      this._resolveExpr(variabledecl.initializer)
    }
    this._define(variabledecl.name)
  }

  visitVariableExprExpr(variableexpr: VariableExpr): void {
    if (
      this._scopes.length &&
      this._scopes[this._scopes.length - 1].get(variableexpr.name.lexeme) === false
    ) {
      this._reportError(variableexpr.name, 'Cannot read local variable in its own initializer.')
    }
    this._resolveLocal(variableexpr, variableexpr.name)
  }

  visitAssignExpr(assign: Assign): void {
    /*
     * First, we resolve the expression for the assigned value
     * in case it also contains references to other variables.
     */
    this._resolveExpr(assign.value)
    this._resolveLocal(assign, assign.name)
  }

  visitFuncStmt(func: Func): void {
    this._declare(func.name)
    this._define(func.name)
    this._resolveFunction(func, FunctionType.Function)
  }

  visitExpressionStmt(expression: Expression): void {
    this._resolveExpr(expression.expression)
  }

  visitIfStmtStmt(ifstmt: IfStmt): void {
    this._resolveExpr(ifstmt.condition)
    this._resolveStmt(ifstmt.thenBranch)
    if (ifstmt.elseBranch) {
      this._resolveStmt(ifstmt.elseBranch)
    }
  }

  visitPrintStmt(print: Print): void {
    this._resolveExpr(print.expression)
  }

  visitReturnStmtStmt(returnstmt: ReturnStmt): void {
    if (this._currentFunction === FunctionType.None) {
      this._reportError(returnstmt.keyword, 'Cannot return from top-level code.')
    }
    if (returnstmt.value) {
      if (this._currentFunction === FunctionType.Initializer) {
        this._reportError(returnstmt.keyword, 'Cannot return a value from an initializer.')
      }
      this._resolveExpr(returnstmt.value)
    }
  }

  visitWhileStmtStmt(whilestmt: WhileStmt): void {
    this._resolveExpr(whilestmt.condition)
    this._resolveStmt(whilestmt.body)
  }

  visitBinaryExpr(binary: Binary): void {
    this._resolveExpr(binary.left)
    this._resolveExpr(binary.right)
  }

  visitCallExpr(call: Call): void {
    this._resolveExpr(call.callee)
    call.args.forEach(arg => this._resolveExpr(arg))
  }

  visitGetExpr(get: Get): void {
    this._resolveExpr(get.obj)
  }

  visitGroupingExpr(grouping: Grouping): void {
    this._resolveExpr(grouping.expression)
  }

  visitLiteralExpr(literal: Literal): void {
    noop()
  }

  visitLogicalExpr(logical: Logical): void {
    this._resolveExpr(logical.left)
    this._resolveExpr(logical.right)
  }

  visitSetExprExpr(setexpr: SetExpr): void {
    this._resolveExpr(setexpr.value)
    this._resolveExpr(setexpr.obj)
  }

  visitSuperExprExpr(superexpr: SuperExpr): void {
    if (this._currentClass === ClassType.None) {
      this._reportError(superexpr.keyword, "Can't use 'super' outside of a class.")
    } else if (this._currentClass !== ClassType.SubClass) {
      this._reportError(superexpr.keyword, "Can't use 'super' in a class with no superclass.")
    }
    this._resolveLocal(superexpr, superexpr.keyword)
  }

  visitThisExprExpr(thisExpr: ThisExpr): void {
    if (this._currentClass === ClassType.None) {
      this._reportError(thisExpr.keyword, "Cannot use 'this' outside of a class.")
      return
    }
    this._resolveLocal(thisExpr, thisExpr.keyword)
  }

  visitUnaryExpr(unary: Unary): void {
    this._resolveExpr(unary.right)
  }

  resolve(statements: Stmt[]): void {
    for (const statement of statements) {
      this._resolveStmt(statement)
    }
  }

  private _beginScope(): void {
    this._scopes.push(new Map())
  }

  private _endScope(): void {
    this._scopes.pop()
  }

  private _resolveStmt(stmt: Stmt): void {
    stmt.accept(this)
  }

  private _resolveExpr(expr: Expr): void {
    expr.accept(this)
  }

  private _declare(name: IToken): void {
    if (!this._scopes.length) return
    const scope = this._scopes[this._scopes.length - 1]
    if (scope.has(name.lexeme)) {
      this._reportError(name, 'Variable with this name already declared in this scope.')
    }
    scope.set(name.lexeme, false)
  }

  private _define(name: IToken): void {
    if (!this._scopes.length) return
    const scope = this._scopes[this._scopes.length - 1]
    scope.set(name.lexeme, true)
  }

  private _resolveLocal(expr: Expr, name: IToken): void {
    for (let i = this._scopes.length - 1; i >= 0; i--) {
      if (this._scopes[i].has(name.lexeme)) {
        this._interpreter.resolve(expr, this._scopes.length - 1 - i)
        return
      }
    }
  }

  private _resolveFunction(func: Func, type: FunctionType): void {
    const enclosingFunction = this._currentFunction
    this._currentFunction = type
    this._beginScope()
    func.params.forEach(param => {
      this._declare(param)
      this._define(param)
    })
    this.resolve(func.body)
    this._endScope()
    this._currentFunction = enclosingFunction
  }
}
