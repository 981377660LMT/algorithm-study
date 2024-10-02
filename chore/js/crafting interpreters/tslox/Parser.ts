/* eslint-disable max-len */

/**
 * Lox grammar.
 *
 * @see {@link https://craftinginterpreters.com/appendix-i.html}
 *
 * program       → declaration* EOF ;
 *
 * declaration   → classDecl | funDecl | varDecl | statement ;
 * classDecl     → "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}" ;
 * funDecl       → "fun" function ;
 * varDecl       → "var" IDENTIFIER ( "=" expression )? ";" ;
 *
 * statement     → exprStatement | forStatement | ifStmt | printStatement | returnStatement | whileStatement | block ;
 * exprStatement → expression ";" ;
 * forStatement  → "for" "(" ( varDecl | exprStatement | ";" ) expression? ";" expression? ")" statement ;
 * ifStmt        → "if" "(" expression ")" statement ( "else" statement )? ;
 * printStatement → "print" expression ";" ;
 * returnStatement → "return" expression? ";" ;
 * whileStatement→ "while" "(" expression ")" statement ;
 * block         → "{" declaration* "}" ;
 *
 * expression     → assignment ;
 * assignment     → ( call "." )? IDENTIFIER "=" assignment | logic_or ;
 * logic_or       → logic_and ( "or" logic_and )* ;
 * logic_and      → equality ( "and" equality )* ;
 * equality       → comparison ( ( "!=" | "==" ) comparison )* ;
 * comparison    → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
 * term          → factor ( ( "-" | "+" ) factor )* ;
 * factor        → unary ( ( "/" | "*" ) unary )* ;
 * unary         → ( "!" | "-" ) unary | call ;
 * call          → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
 * primary       → "true" | "false" | "nil" | "this" | NUMBER | STRING | IDENTIFIER | "(" expression ")" | "super" "." IDENTIFIER ;
 *
 * function     → IDENTIFIER "(" parameters? ")" block ;
 * parameters   → IDENTIFIER ( "," IDENTIFIER )* ;
 * arguments    → expression ( "," expression )* ;
 */

import { ParseError } from './consts'
import {
  Expr,
  Binary,
  Grouping,
  Literal,
  Unary,
  Stmt,
  Print,
  Expression,
  VariableDecl,
  VariableExpr,
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
import { IToken, ReportErrorFunc, TokenType } from './types'

export class Parser {
  private readonly _tokens: IToken[]
  private _current = 0

  private _reportError: ReportErrorFunc

  constructor(tokens: IToken[], options: { reportError?: ReportErrorFunc } = {}) {
    this._tokens = tokens
    // eslint-disable-next-line no-console
    this._reportError = options.reportError || console.error
  }

  parse(): Stmt[] {
    const res: Stmt[] = []
    while (!this._isAtEnd()) {
      const v = this._declaration()
      if (v) res.push(v)
    }
    return res
  }

  /**
   * First, it looks to see if we’re at a variable declaration by looking for the leading var keyword.
   * If not, it falls back to parsing a statement.
   */
  private _declaration(): Stmt | undefined {
    try {
      if (this._match(TokenType.CLASS)) return this._classDeclaration()
      if (this._match(TokenType.FUN)) return this._func('function')
      if (this._match(TokenType.VAR)) return this._varDeclaration()
      return this._statement()
    } catch (error) {
      if (error instanceof ParseError) {
        this._synchronize()
      }
      return undefined
    }
  }

  private _classDeclaration(): Stmt {
    const name = this._consume(TokenType.IDENTIFIER, 'Expect class name.')
    let superclass: VariableExpr | undefined
    if (this._match(TokenType.LESS)) {
      this._consume(TokenType.IDENTIFIER, 'Expect superclass name.')
      superclass = new VariableExpr(this._previous())
    }
    this._consume(TokenType.LEFT_BRACE, 'Expect "{" before class body.')

    const methods: Func[] = []
    while (!this._check(TokenType.RIGHT_BRACE) && !this._isAtEnd()) {
      methods.push(this._func('method'))
    }

    this._consume(TokenType.RIGHT_BRACE, 'Expect "}" after class body.')
    return new ClassStmt(name, superclass, methods)
  }

  private _varDeclaration(): Stmt {
    const name = this._consume(TokenType.IDENTIFIER, 'Expect variable name.')
    let initializer: Expr | undefined
    if (this._match(TokenType.EQUAL)) {
      initializer = this._expression()
    }
    this._consume(TokenType.SEMICOLON, 'Expect ";" after variable declaration.')
    return new VariableDecl(name, initializer)
  }

  private _statement(): Stmt {
    if (this._match(TokenType.FOR)) return this._forStatement()
    if (this._match(TokenType.IF)) return this._ifStatement()
    if (this._match(TokenType.PRINT)) return this._printStatement()
    if (this._match(TokenType.RETURN)) return this._returnStatement()
    if (this._match(TokenType.WHILE)) return this._whileStatement()
    if (this._match(TokenType.LEFT_BRACE)) return new Block(this._block())
    return this._expressionStatement()
  }

  private _forStatement(): Stmt {
    this._consume(TokenType.LEFT_PAREN, 'Expect "(" after "for".')

    let initializer: Stmt | undefined
    if (this._match(TokenType.SEMICOLON)) {
      initializer = undefined
    } else if (this._match(TokenType.VAR)) {
      initializer = this._varDeclaration()
    } else {
      initializer = this._expressionStatement()
    }

    let condition: Expr | undefined
    if (!this._check(TokenType.SEMICOLON)) {
      condition = this._expression()
    }
    this._consume(TokenType.SEMICOLON, 'Expect ";" after loop condition.')

    let increment: Expr | undefined
    if (!this._check(TokenType.RIGHT_PAREN)) {
      increment = this._expression()
    }
    this._consume(TokenType.RIGHT_PAREN, 'Expect ")" after for clauses.')

    let body = this._statement()
    if (increment) body = new Block([body, new Expression(increment)])
    if (!condition) condition = new Literal(true)
    body = new WhileStmt(condition, body)
    if (initializer) body = new Block([initializer, body])

    return body
  }

  private _ifStatement(): IfStmt {
    this._consume(TokenType.LEFT_PAREN, 'Expect "(" after "if".')
    const condition = this._expression()
    this._consume(TokenType.RIGHT_PAREN, 'Expect ")" after if condition.')

    const thenBranch = this._statement()
    let elseBranch: Stmt | undefined
    if (this._match(TokenType.ELSE)) {
      elseBranch = this._statement()
    }

    return new IfStmt(condition, thenBranch, elseBranch)
  }

  private _printStatement(): Print {
    const value = this._expression()
    this._consume(TokenType.SEMICOLON, 'Expect ";" after value.')
    return new Print(value)
  }

  private _returnStatement(): ReturnStmt {
    const keyword = this._previous()
    let value: Expr | undefined
    if (!this._check(TokenType.SEMICOLON)) {
      value = this._expression()
    }
    this._consume(TokenType.SEMICOLON, 'Expect ";" after return value.')
    return new ReturnStmt(keyword, value)
  }

  private _whileStatement(): WhileStmt {
    this._consume(TokenType.LEFT_PAREN, 'Expect "(" after "while".')
    const condition = this._expression()
    this._consume(TokenType.RIGHT_PAREN, 'Expect ")" after condition.')
    const body = this._statement()
    return new WhileStmt(condition, body)
  }

  /** block() assumes the brace token has already been matched. */
  private _block(): Stmt[] {
    const res: Stmt[] = []
    while (!this._check(TokenType.RIGHT_BRACE) && !this._isAtEnd()) {
      const v = this._declaration()
      if (v) res.push(v)
    }
    this._consume(TokenType.RIGHT_BRACE, 'Expect "}" after block.')
    return res
  }

  private _expressionStatement(): Expression {
    const value = this._expression()
    this._consume(TokenType.SEMICOLON, 'Expect ";" after value.')
    return new Expression(value)
  }

  private _func(kind: 'function' | 'method'): Func {
    const name = this._consume(TokenType.IDENTIFIER, `Expect ${kind} name.`)
    this._consume(TokenType.LEFT_PAREN, `Expect "(" after ${kind} name.`)
    const params: IToken[] = []
    if (!this._check(TokenType.RIGHT_PAREN)) {
      while (true) {
        if (params.length >= 255) {
          this._error(this._peek(), 'Cannot have more than 255 parameters.')
        }
        params.push(this._consume(TokenType.IDENTIFIER, 'Expect parameter name.'))
        if (!this._match(TokenType.COMMA)) break
      }
    }
    this._consume(TokenType.RIGHT_PAREN, 'Expect ")" after parameters.')
    this._consume(TokenType.LEFT_BRACE, `Expect "{" before ${kind} body.`)
    const body = this._block()
    return new Func(name, params, body)
  }

  private _expression(): Expr {
    return this._assignment()
  }

  private _assignment(): Expr {
    const expr = this._or()
    /** We parse the left-hand side, which can be any expression of higher precedence. */
    if (this._match(TokenType.EQUAL)) {
      const equals = this._previous()
      const value = this._assignment()
      if (expr instanceof VariableExpr) {
        return new Assign(expr.name, value)
      }
      if (expr instanceof Get) {
        return new SetExpr(expr.obj, expr.name, value)
      }
      this._error(equals, 'Invalid assignment target.')
    }
    return expr
  }

  private _or(): Expr {
    let expr = this._and()
    while (this._match(TokenType.OR)) {
      const operator = this._previous()
      const right = this._and()
      expr = new Logical(expr, operator, right)
    }
    return expr
  }

  private _and(): Expr {
    let expr = this._equality()
    while (this._match(TokenType.AND)) {
      const operator = this._previous()
      const right = this._equality()
      expr = new Logical(expr, operator, right)
    }
    return expr
  }

  private _equality(): Expr {
    let res = this._comparison()
    while (this._match(TokenType.BANG_EQUAL, TokenType.EQUAL_EQUAL)) {
      const operator = this._previous()
      const right = this._comparison()
      res = new Binary(res, operator, right)
    }
    return res
  }

  private _comparison(): Expr {
    let res = this._term()
    while (
      this._match(TokenType.GREATER, TokenType.GREATER_EQUAL, TokenType.LESS, TokenType.LESS_EQUAL)
    ) {
      const operator = this._previous()
      const right = this._term()
      res = new Binary(res, operator, right)
    }
    return res
  }

  private _term(): Expr {
    let res = this._factor()
    while (this._match(TokenType.MINUS, TokenType.PLUS)) {
      const operator = this._previous()
      const right = this._factor()
      res = new Binary(res, operator, right)
    }
    return res
  }

  private _factor(): Expr {
    let res = this._unary()
    while (this._match(TokenType.SLASH, TokenType.STAR)) {
      const operator = this._previous()
      const right = this._unary()
      res = new Binary(res, operator, right)
    }
    return res
  }

  private _unary(): Expr {
    if (this._match(TokenType.BANG, TokenType.MINUS)) {
      const operator = this._previous()
      const right = this._unary()
      return new Unary(operator, right)
    }
    return this._call()
  }

  private _call(): Expr {
    let expr = this._primary()
    while (true) {
      if (this._match(TokenType.LEFT_PAREN)) {
        expr = this._finishCall(expr)
      } else if (this._match(TokenType.DOT)) {
        const name = this._consume(TokenType.IDENTIFIER, 'Expect property name after ".".')
        expr = new Get(expr, name)
      } else {
        break
      }
    }
    return expr
  }

  private _finishCall(callee: Expr): Expr {
    const args: Expr[] = []
    // no arguments
    if (!this._check(TokenType.RIGHT_PAREN)) {
      while (true) {
        // reports the error and keeps on keepin’ on.
        if (args.length >= 255) {
          this._error(this._peek(), 'Cannot have more than 255 arguments.')
        }
        args.push(this._expression())
        if (!this._match(TokenType.COMMA)) break
      }
    }
    const paren = this._consume(TokenType.RIGHT_PAREN, 'Expect ")" after arguments.')
    return new Call(callee, paren, args)
  }

  private _primary(): Expr {
    if (this._match(TokenType.FALSE)) return new Literal(false)
    if (this._match(TokenType.TRUE)) return new Literal(true)
    if (this._match(TokenType.NIL)) return new Literal(undefined)
    if (this._match(TokenType.NUMBER, TokenType.STRING)) {
      return new Literal(this._previous().literal)
    }
    if (this._match(TokenType.SUPER)) {
      const keyword = this._previous()
      this._consume(TokenType.DOT, 'Expect "." after "super".')
      const method = this._consume(TokenType.IDENTIFIER, 'Expect superclass method name.')
      return new SuperExpr(keyword, method)
    }
    if (this._match(TokenType.THIS)) return new ThisExpr(this._previous())

    if (this._match(TokenType.IDENTIFIER)) {
      return new VariableExpr(this._previous())
    }

    if (this._match(TokenType.LEFT_PAREN)) {
      const expr = this._expression()
      this._consume(TokenType.RIGHT_PAREN, 'Expect ")" after expression.')
      return new Grouping(expr)
    }

    // !None of the cases in there match
    throw this._error(this._peek(), 'Expect expression.')
  }

  /** Consumes the current token if it matches the given type, otherwise throws an error. */
  private _consume(type: TokenType, message: string): IToken {
    if (this._check(type)) return this._advance()
    throw this._error(this._peek(), message)
  }

  private _match(...types: TokenType[]): boolean {
    for (let i = 0; i < types.length; i++) {
      if (this._check(types[i])) {
        this._advance()
        return true
      }
    }
    return false
  }

  private _check(type: TokenType): boolean {
    if (this._isAtEnd()) return false
    return this._peek().type === type
  }

  /** Consumes the current token and returns it. */
  private _advance(): IToken {
    if (!this._isAtEnd()) this._current++
    return this._previous()
  }

  private _isAtEnd(): boolean {
    return this._peek().type === TokenType.EOF
  }

  /** Returns the current token we have yet to consume. */
  private _peek(): IToken {
    return this._tokens[this._current]
  }

  /** Returns the most recently consumed token. */
  private _previous(): IToken {
    return this._tokens[this._current - 1]
  }

  private _error(token: IToken, message: string): ParseError {
    this._reportError(token, message)
    return new ParseError()
  }

  private _synchronize(): void {
    this._advance()
    while (!this._isAtEnd()) {
      if (this._previous().type === TokenType.SEMICOLON) return
      switch (this._peek().type) {
        case TokenType.CLASS:
        case TokenType.FUN:
        case TokenType.VAR:
        case TokenType.FOR:
        case TokenType.IF:
        case TokenType.WHILE:
        case TokenType.PRINT:
        case TokenType.RETURN:
          return
        default:
          break
      }
      this._advance()
    }
  }
}
