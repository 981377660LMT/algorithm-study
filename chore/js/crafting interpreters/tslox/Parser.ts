/* eslint-disable max-len */

/**
 * program       → declaration* EOF ;
 * declaration   → varDecl | statement ;
 * varDecl       → "var" IDENTIFIER ( "=" expression )? ";" ;
 * statement     → exprStatement | printStatement ;
 * exprStatement → expression ";" ;
 * printStatement → "print" expression ";" ;
 *
 * expression     → equality ;
 * equality       → comparison ( ( "!=" | "==" ) comparison )* ;
 * comparison    → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
 * term          → factor ( ( "-" | "+" ) factor )* ;
 * factor        → unary ( ( "/" | "*" ) unary )* ;
 * unary         → ( "!" | "-" ) unary | primary ;
 * primary       → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
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
  VariableExpr
} from './expr'
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

  parse(): (Stmt | undefined)[] {
    const res: (Stmt | undefined)[] = []
    while (!this._isAtEnd()) {
      res.push(this._declaration())
    }
    return res
  }

  /**
   * First, it looks to see if we’re at a variable declaration by looking for the leading var keyword.
   * If not, it falls back to parsing a statement.
   */
  private _declaration(): Stmt | undefined {
    try {
      if (this._match(TokenType.VAR)) return this._varDeclaration()
      return this._statement()
    } catch (error) {
      if (error instanceof ParseError) {
        this._synchronize()
      }
      return undefined
    }
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
    if (this._match(TokenType.PRINT)) return this._printStatement()
    return this._expressionStatement()
  }

  private _printStatement(): Stmt {
    const value = this._expression()
    this._consume(TokenType.SEMICOLON, 'Expect ";" after value.')
    return new Print(value)
  }

  private _expressionStatement(): Stmt {
    const value = this._expression()
    this._consume(TokenType.SEMICOLON, 'Expect ";" after value.')
    return new Expression(value)
  }

  private _expression(): Expr {
    return this._equality()
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
    return this._primary()
  }

  private _primary(): Expr {
    if (this._match(TokenType.FALSE)) return new Literal(false)
    if (this._match(TokenType.TRUE)) return new Literal(true)
    if (this._match(TokenType.NIL)) return new Literal(undefined)
    if (this._match(TokenType.NUMBER, TokenType.STRING)) {
      return new Literal(this._previous().literal)
    }

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
