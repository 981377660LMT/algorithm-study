/* eslint-disable class-methods-use-this */

import { RuntimeError } from '../consts'
import { type Interpreter } from '../Interpreter'
import { IToken } from '../types'
import { LoxCallable } from './LoxCallable'
import { type LoxFunction } from './LoxFunction'

/**
 * For Behavior.
 */
export class LoxClass extends LoxCallable {
  readonly name: string
  readonly superClass: LoxClass | undefined
  private readonly _methods: Map<string, LoxFunction>

  constructor(name: string, superClass: LoxClass | undefined, methods: Map<string, LoxFunction>) {
    super()
    this.name = name
    this.superClass = superClass
    this._methods = methods
  }

  call(interpter: Interpreter, args: unknown[]): unknown {
    const instance = new LoxInstance(this)
    const initializer = this.findMethod('init')
    if (initializer) initializer.bind(instance).call(interpter, args)
    return instance
  }

  arity(): number {
    const initializer = this.findMethod('init')
    if (initializer) return initializer.arity()
    return 0
  }

  /**
   * Find method along the inheritance chain.
   */
  findMethod(name: string): LoxFunction | undefined {
    if (this._methods.has(name)) return this._methods.get(name)
    if (this.superClass) return this.superClass.findMethod(name)
    return undefined
  }

  override toString(): string {
    return this.name
  }
}

/**
 * For State.
 */
export class LoxInstance {
  private readonly _cls: LoxClass
  private readonly _fields = new Map<string, unknown>()

  constructor(cls: LoxClass) {
    this._cls = cls
  }

  get(name: IToken): unknown {
    if (this._fields.has(name.lexeme)) {
      return this._fields.get(name.lexeme)
    }

    const method = this._cls.findMethod(name.lexeme)
    if (method) return method.bind(this)

    throw new RuntimeError(name, `Undefined property '${name.lexeme}'.`)
  }

  set(name: IToken, value: unknown): void {
    this._fields.set(name.lexeme, value)
  }

  toString(): string {
    return `${this._cls.name} instance`
  }
}
