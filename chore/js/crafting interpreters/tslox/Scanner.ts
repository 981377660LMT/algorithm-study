/* eslint-disable no-lone-blocks */

import { createToken } from './Token'
import { type IScanner, type IToken, type ReportErrorFunc, TokenType } from './types'
import { KEY_WORDS } from './consts'
import { isAlpha, isDigit, isAlphaNumeric } from './utils'

export class Scanner implements IScanner<IToken> {
  private readonly _source: string
  private readonly _tokens: IToken[] = []

  /** The beginning of the current lexeme being scanned. */
  private _start = 0
  /** The current character being scanned. */
  private _current = 0
  /** The current line being scanned. */
  private _line = 1

  private _reportError: ReportErrorFunc

  constructor(source: string, options: { reportError?: ReportErrorFunc } = {}) {
    this._source = source
    // eslint-disable-next-line no-console
    this._reportError = options.reportError || console.error
  }

  scanTokens(): IToken[] {
    while (!this._isAtEnd()) {
      // We are at the beginning of the next lexeme.
      this._start = this._current
      this._scanToken()
    }

    this._tokens.push(
      createToken({ type: TokenType.EOF, lexeme: '', literal: undefined, line: this._line })
    )
    return this._tokens
  }

  private _isAtEnd(): boolean {
    return this._current >= this._source.length
  }

  private _scanToken(): void {
    const c = this._advance()

    switch (c) {
      // Single-character tokens.
      case '(':
        this._addToken(TokenType.LEFT_PAREN)
        break
      case ')':
        this._addToken(TokenType.RIGHT_PAREN)
        break
      case '{':
        this._addToken(TokenType.LEFT_BRACE)
        break
      case '}':
        this._addToken(TokenType.RIGHT_BRACE)
        break
      case ',':
        this._addToken(TokenType.COMMA)
        break
      case '.':
        this._addToken(TokenType.DOT)
        break
      case '-':
        this._addToken(TokenType.MINUS)
        break
      case '+':
        this._addToken(TokenType.PLUS)
        break
      case ';':
        this._addToken(TokenType.SEMICOLON)
        break
      case '*':
        this._addToken(TokenType.STAR)
        break
      // division or comment
      case '/':
        if (this._match('/')) {
          // A comment goes until the end of the line.
          while (this._peek() !== '\n' && !this._isAtEnd()) this._advance()
        } else {
          this._addToken(TokenType.SLASH)
        }
        break

      // One or two character tokens.
      case '!':
        this._addToken(this._match('=') ? TokenType.BANG_EQUAL : TokenType.BANG)
        break
      case '=':
        this._addToken(this._match('=') ? TokenType.EQUAL_EQUAL : TokenType.EQUAL)
        break
      case '<':
        this._addToken(this._match('=') ? TokenType.LESS_EQUAL : TokenType.LESS)
        break
      case '>':
        this._addToken(this._match('=') ? TokenType.GREATER_EQUAL : TokenType.GREATER)
        break

      // Ignore whitespace.
      case ' ':
      case '\r':
      case '\t':
        break

      case '\n':
        this._line++
        break

      case '"':
        this._string()
        break

      default:
        if (isDigit(c)) {
          this._number()
        } else if (isAlpha(c)) {
          this._identifier()
        } else {
          this._reportError(this._line, `Unexpected character: ${c}`)
        }

        break
    }
  }

  /** Consumes the next character in the source file and returns it. */
  private _advance(): string {
    this._current++
    return this._source[this._current - 1]
  }

  private _addToken(type: TokenType, literal: unknown = undefined): void {
    const text = this._source.slice(this._start, this._current)
    this._tokens.push(createToken({ type, literal, lexeme: text, line: this._line }))
  }

  private _match(expected: string): boolean {
    if (this._isAtEnd()) return false
    if (this._source[this._current] !== expected) return false
    this._current++
    return true
  }

  /** lookahead */
  private _peek(): string {
    if (this._isAtEnd()) return '\0'
    return this._source[this._current]
  }

  private _peekNext(): string {
    if (this._current + 1 >= this._source.length) return '\0'
    return this._source[this._current + 1]
  }

  private _string(): void {
    while (this._peek() !== '"' && !this._isAtEnd()) {
      if (this._peek() === '\n') this._line++
      this._advance()
    }

    if (this._isAtEnd()) {
      this._reportError(this._line, 'Unterminated string')
      return
    }

    // The closing ".
    this._advance()

    // Trim the surrounding quotes.
    const value = this._source.slice(this._start + 1, this._current - 1)
    this._addToken(TokenType.STRING, value)
  }

  private _number(): void {
    while (isDigit(this._peek())) this._advance()
    // Look for a fractional part.
    if (this._peek() === '.' && isDigit(this._peekNext())) {
      // Consume the "."
      this._advance()
      while (isDigit(this._peek())) this._advance()
    }

    this._addToken(
      TokenType.NUMBER,
      Number.parseFloat(this._source.slice(this._start, this._current))
    )
  }

  private _identifier(): void {
    while (isAlphaNumeric(this._peek())) this._advance()
    const text = this._source.slice(this._start, this._current)
    const type = KEY_WORDS.get(text) || TokenType.IDENTIFIER
    this._addToken(type)
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  {
    const scanne = new Scanner('(){}')
    console.log(scanne.scanTokens())
  }

  {
    const scanner = new Scanner('(){}//')
    console.log(scanner.scanTokens())
  }

  {
    const scanner = new Scanner('!*+-/=<> <= == ')
    console.log(scanner.scanTokens())
  }

  {
    const scanner = new Scanner('"hello world"')
    console.log(scanner.scanTokens())
  }

  {
    const scanner = new Scanner('123 123.456')
    console.log(scanner.scanTokens())
  }

  {
    const scanner = new Scanner(
      'and class else false for fun if nil or print return super this true var while'
    )
    console.log(scanner.scanTokens())
  }

  {
    const scanner = new Scanner('var a = 1;')
    console.log(scanner.scanTokens())
  }
}
