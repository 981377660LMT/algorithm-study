import { Return } from '../consts'
import { Environment } from '../Environment'
import { type Interpreter } from '../Interpreter'
import { Func } from '../syntaxTree'
import { LoxCallable } from './LoxCallable'
import { type LoxInstance } from './LoxClass'

export class LoxFunction extends LoxCallable {
  private _declaration!: Func
  /** When we create a LoxFunction, we capture the current environment. */
  private _closure!: Environment
  private _isInitializer!: boolean
  private _inited = false

  init(declaration: Func, closure: Environment, isInitializer = false): LoxFunction {
    this._declaration = declaration
    this._closure = closure
    this._isInitializer = isInitializer
    this._inited = true
    return this
  }

  /**
   * Each function call gets its own environment.
   */
  call(interpreter: Interpreter, args: unknown[]): unknown {
    if (!this._inited) {
      throw new Error('LoxFunction not inited.')
    }
    const env = new Environment(this._closure)
    this._declaration.params.forEach((p, i) => {
      env.define(p.lexeme, args[i])
    })

    try {
      interpreter.executeBlock(this._declaration.body, env)
    } catch (error) {
      if (error instanceof Return) {
        if (this._isInitializer) return this._closure.getAt(0, 'this')
        return error.value
      }
      throw error
    }

    /** init() methods always return this. */
    if (this._isInitializer) return this._closure.getAt(0, 'this')

    return undefined
  }

  bind(instance: LoxInstance): LoxFunction {
    const env = new Environment(this._closure)
    env.define('this', instance)
    return new LoxFunction().init(this._declaration, env, this._isInitializer)
  }

  arity(): number {
    if (!this._inited) {
      throw new Error('LoxFunction not inited.')
    }
    return this._declaration.params.length
  }

  override toString(): string {
    return `<fn ${this._declaration.name.lexeme}>`
  }
}
