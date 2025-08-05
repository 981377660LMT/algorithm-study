interface Rollbackable<T> {
  snapShot(): T
  rollback(to: T): void
}

class RValue<V> implements Rollbackable<V> {
  private _value: V

  constructor(value: V) {
    this._value = value
  }

  set(value: V): void {
    this._value = value
  }

  get(): V {
    return this._value
  }

  snapShot(): V {
    return this._value
  }

  rollback(to: V): void {
    this._value = to
  }
}

class RArray<V> implements Rollbackable<number> {
  private readonly _arr: V[]
  private readonly _history: { index: number; value: V }[]

  constructor(arr: V[]) {
    this._arr = arr.slice()
    this._history = []
  }

  get(index: number): V {
    return this._arr[index]
  }

  set(index: number, value: V): void {
    this._history.push({ index, value: this._arr[index] })
    this._arr[index] = value
  }

  snapShot(): number {
    return this._history.length
  }

  rollback(to: number): void {
    for (let i = this._history.length - 1; i >= to; i--) {
      const { index, value } = this._history[i]
      this._arr[index] = value
    }
    this._history.length = to
  }
}

export {}
