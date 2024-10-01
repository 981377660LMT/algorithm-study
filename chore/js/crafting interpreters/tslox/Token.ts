import { IToken, TokenType } from './types'

export const createToken = ({
  type,
  lexeme,
  literal,
  line
}: {
  type: TokenType
  lexeme: string
  literal: unknown
  line: number
}): IToken => new Token(type, lexeme, literal, line)

class Token implements IToken {
  readonly type: TokenType
  readonly lexeme: string
  readonly literal: unknown
  readonly line: number

  constructor(type: TokenType, lexeme: string, literal: unknown, line: number) {
    this.type = type
    this.lexeme = lexeme
    this.literal = literal
    this.line = line
  }

  toString(): string {
    return `${this.type} ${this.lexeme} ${this.literal}`
  }
}
