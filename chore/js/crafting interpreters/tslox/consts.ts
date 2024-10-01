import { IToken, TokenType } from './types'

/** Lox 语言保留字. */
export const KEY_WORDS = new Map<string, TokenType>()
KEY_WORDS.set('and', TokenType.AND)
KEY_WORDS.set('class', TokenType.CLASS)
KEY_WORDS.set('else', TokenType.ELSE)
KEY_WORDS.set('false', TokenType.FALSE)
KEY_WORDS.set('for', TokenType.FOR)
KEY_WORDS.set('fun', TokenType.FUN)
KEY_WORDS.set('if', TokenType.IF)
KEY_WORDS.set('nil', TokenType.NIL)
KEY_WORDS.set('or', TokenType.OR)
KEY_WORDS.set('print', TokenType.PRINT)
KEY_WORDS.set('return', TokenType.RETURN)
KEY_WORDS.set('super', TokenType.SUPER)
KEY_WORDS.set('this', TokenType.THIS)
KEY_WORDS.set('true', TokenType.TRUE)
KEY_WORDS.set('var', TokenType.VAR)
KEY_WORDS.set('while', TokenType.WHILE)

export class ParseError extends Error {}

export class RuntimeError extends Error {
  readonly token: IToken

  constructor(token: IToken, message: string) {
    super(message)
    this.token = token
  }
}
