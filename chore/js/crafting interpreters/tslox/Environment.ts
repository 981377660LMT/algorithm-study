import { RuntimeError } from './consts'
import { IToken } from './types'

/**
 * TODO:
 *  A more efficient environment representation would
 *  store local variables in an **array** and look them up by index.
 */
export class Environment {
  /** The enclosing environment. */
  readonly enclosing: Environment | undefined
  private readonly _mp: Map<string, any> = new Map()

  constructor(enclosing?: Environment) {
    this.enclosing = enclosing
  }

  define(name: string, value: any): void {
    this._mp.set(name, value)
  }

  get(name: IToken): any {
    if (this._mp.has(name.lexeme)) {
      return this._mp.get(name.lexeme)
    }
    if (this.enclosing) return this.enclosing.get(name)
    throw new RuntimeError(name, `Undefined variable '${name.lexeme}'.`)
  }

  getAt(distance: number, name: string): any {
    return this._ancestor(distance)._mp.get(name)
  }

  /**
   * The key difference between assignment and definition is that
   * assignment is not allowed to create a new variable.
   */
  assign(name: IToken, value: any): void {
    if (this._mp.has(name.lexeme)) {
      this._mp.set(name.lexeme, value)
      return
    }
    if (this.enclosing) {
      this.enclosing.assign(name, value)
      return
    }
    throw new RuntimeError(name, `Undefined variable '${name.lexeme}'.`)
  }

  assignAt(distance: number, name: IToken, value: any): void {
    this._ancestor(distance)._mp.set(name.lexeme, value)
  }

  private _ancestor(distance: number): Environment {
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    let env: Environment = this
    for (let i = 0; i < distance; i++) env = env.enclosing!
    return env
  }
}
