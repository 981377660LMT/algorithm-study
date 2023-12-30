/**
 * 惰性求值.
 */
class Suspension<S> {
  private readonly _x: S | (() => S) | undefined = undefined
  private _resolved: S | undefined = undefined

  constructor(x?: S | (() => S)) {
    this._x = x
  }

  resolve(): S | undefined {
    // eslint-disable-next-line eqeqeq
    if (this._resolved == undefined) {
      this._resolved = typeof this._x === 'function' ? (this._x as () => S)() : this._x
    }
    return this._resolved
  }
}

export { Suspension }

if (require.main === module) {
  const s = new Suspension(() => Date.now())
  console.log(s.resolve())
}
