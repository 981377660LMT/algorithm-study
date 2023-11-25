/* eslint-disable @typescript-eslint/no-non-null-assertion */

class RollbackArray<T> {
  private readonly _n: number
  private readonly _data: T[]
  private readonly _history: { index: number; value: T }[] = []

  constructor(arr: T[])
  constructor(n: number, f: (index: number) => T)
  constructor(arg1: any, arg2?: any) {
    if (Array.isArray(arg1)) {
      this._n = arg1.length
      this._data = arg1.slice()
    } else {
      this._n = arg1
      this._data = Array(arg1)
      for (let i = 0; i < arg1; i++) {
        this._data[i] = arg2(i)
      }
    }
  }

  getTime(): number {
    return this._history.length
  }

  rollback(time: number): void {
    while (this._history.length > time) {
      const { index, value } = this._history.pop()!
      this._data[index] = value
    }
  }

  undo(): boolean {
    if (!this._history.length) return false
    const { index, value } = this._history.pop()!
    this._data[index] = value
    return true
  }

  get(index: number): T {
    return this._data[index]
  }

  set(index: number, value: T): void {
    this._history.push({ index, value: this._data[index] })
    this._data[index] = value
  }

  getAll(): T[] {
    return this._data.slice()
  }

  get length(): number {
    return this._n
  }
}

export { RollbackArray }

if (require.main === module) {
  const arr = new RollbackArray(10, () => 0)
  arr.set(0, 1)
  console.log(arr.getAll())
  arr.undo()
  console.log(arr.getAll())
}
