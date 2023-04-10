// https://www.cnblogs.com/bzy-blog/p/14353073.html

/**
 * 线性序列并查集(NextFinder).
 */
class LinearSequenceUnionFind {
  readonly right: Uint32Array
  private readonly _n: number
  private readonly _data: Uint32Array

  constructor(n: number) {
    this._n = n
    const len = (n >> 5) + 1
    const right = new Uint32Array(len)
    const data = new Uint32Array(len)
    for (let i = 0; i < len; i++) {
      right[i] = i
      data[i] = -1
    }
    this.right = right
    this._data = data
  }

  next(x: number): number | null {
    if (x < 0) x = 0
    if (x >= this._n) return null
    let div = x >> 5
    let mod = x & 31
    let mask = this._data[div] >> mod
    if (mask) {
      // !trailingZeros32(x): 31 - Math.clz32(x & -x)
      const res = ((div << 5) | mod) + 31 - Math.clz32(mask & -mask)
      return res < this._n ? res : null
    }
    div = this._findNext(div + 1)
    const tmp = this._data[div]
    const res = (div << 5) + 31 - Math.clz32(tmp & -tmp)
    return res < this._n ? res : null
  }

  erase(x: number): void {
    let div = x >> 5
    let mod = x & 31
    // flip
    if ((this._data[div] >> mod) & 1) {
      this._data[div] ^= 1 << mod
    }
    if (!this._data[div]) {
      this.right[div] = div + 1 // union to right
    }
  }

  has(x: number): boolean {
    return !!((this._data[x >> 5] >> (x & 31)) & 1)
  }

  toString(): string {
    const res: number[] = []
    for (let i = 0; i < this._n; i++) {
      if (this.has(i)) res.push(i)
    }
    return `Finder(${res.join(',')})`
  }

  private _findNext(x: number): number {
    if (this.right[x] === x) return x
    this.right[x] = this._findNext(this.right[x])
    return this.right[x]
  }
}

if (require.main === module) {
  const uf = new LinearSequenceUnionFind(10)
  uf.erase(0)

  console.log(uf.next(0))
  console.log(uf.next(2))
  console.log(uf.has(0))
  uf.erase(2)

  console.log(uf.next(2))
  console.log(uf.next(9))
  uf.erase(9)
  console.log(uf.next(9))

  console.log(uf.toString())
}

export { LinearSequenceUnionFind, LinearSequenceUnionFind as Finder }
