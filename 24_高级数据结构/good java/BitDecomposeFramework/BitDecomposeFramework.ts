class BitDecomposeFramework<E> {
  private readonly _bits: (E | undefined)[]
  private readonly _merger: (a: E, b: E) => E
  private _end = 0

  constructor(merger: (a: E, b: E) => E) {
    this._merger = merger
    this._bits = Array(32).fill(undefined) // int32
  }

  add(value: E): void {
    this._add(value, 0)
  }

  addAll(other: BitDecomposeFramework<E>): void {
    for (let i = 0; i < other._end; i++) {
      if (other._bits[i] != undefined) {
        this._add(other._bits[i] as E, i)
      }
    }
  }

  forEach(callback: (value: E, index: number) => boolean | void): void {
    for (let i = 0; i < this._end; i++) {
      if (this._bits[i] != undefined) {
        if (callback(this._bits[i] as E, i)) {
          return
        }
      }
    }
  }

  private _add(value: E, index: number): void {
    if (this._bits[index] == undefined) {
      this._bits[index] = value
      this._end = Math.max(this._end, index + 1)
      return
    }
    this._add(this._merger(this._bits[index] as E, value), index + 1)
    this._bits[index] = undefined
  }
}

export { BitDecomposeFramework }

if (require.main === module) {
  const f1 = new BitDecomposeFramework<number>((a, b) => a + b)
  f1.add(1)
  f1.add(2)
  f1.add(3)
  f1.add(4)
  f1.add(5)
  f1.add(6)
  f1.add(16)
  f1.forEach(console.log)
}
