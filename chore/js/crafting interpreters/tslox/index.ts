/* eslint-disable no-lone-blocks */
/* eslint-disable no-console */

// TODO:
// 1. Dependens on interface, not class.
// 2. Service locator, not only components.

import { Scanner } from './Scanner'
import { Parser } from './Parser'
import { Interpreter } from './Interpreter'
import { IToken, TokenType } from './types'
import { RuntimeError } from './consts'

export class TsLox {
  private _hadError = false
  private _hadRuntimeError = false

  /**
   * 在一个REPL会话中连续调用run()时重复使用同一个解释器.
   * 目前这一点没有什么区别，但以后当解释器需要存储全局变量时就会有区别。
   * 这些全局变量应该在整个REPL会话中持续存在。
   */
  private readonly _interpreter: Interpreter

  constructor() {
    this._interpreter = new Interpreter({ reportError: this.runtimeError.bind(this) })
  }

  run(source: string): void {
    const scanner = new Scanner(source, { reportError: this.error.bind(this) })
    const tokens = scanner.scanTokens()
    const parser = new Parser(tokens, { reportError: this.error.bind(this) })
    const statements = parser.parse()
    if (this._hadError || this._hadRuntimeError) return
    if (!statements) return
    this._interpreter.interpret(statements)
  }

  error(pos: number | IToken, message: string): void {
    if (typeof pos === 'number') {
      this._report(pos, '', message)
    } else if (pos.type === TokenType.EOF) {
      this._report(pos.line, 'at end', message)
    } else {
      this._report(pos.line, `at '${pos.lexeme}'`, message)
    }
  }

  runtimeError(error: RuntimeError): void {
    console.log(`${error.message}\n[line ${error.token.line}]`)
    this._hadRuntimeError = true
  }

  private _report(line: number, where: string, message: string): void {
    console.log(`[line ${line}] Error ${where}: ${message}`)
    this._hadError = true
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  {
    const lox = new TsLox()
    lox.run('print (1 + 2) / 3;')
  }
}
