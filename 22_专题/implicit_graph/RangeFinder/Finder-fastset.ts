// 寻找前驱后继/区间删除
// 非常快

/**
 * 利用位运算寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 * 初始时,所有位置都未被访问过.
 */
class Finder {
  private readonly _n: number
  private readonly _lg: number
  private readonly _seg: Uint32Array[]

  constructor(n: number) {
    this._n = n
    const seg: Uint32Array[] = []
    while (true) {
      seg.push(new Uint32Array((n + 31) >>> 5))
      n = (n + 31) >>> 5
      if (n <= 1) {
        break
      }
    }
    this._lg = seg.length
    this._seg = seg
    for (let i = 0; i < this._n; i++) {
      this.insert(i)
    }
  }

  insert(i: number): void {
    for (let h = 0; h < this._lg; h++) {
      this._seg[h][i >>> 5] |= 1 << (i & 31)
      i >>>= 5
    }
  }

  has(i: number): boolean {
    return !!(this._seg[0][i >>> 5] & (1 << (i & 31)))
  }

  erase(i: number): void {
    for (let h = 0; h < this._lg; h++) {
      this._seg[h][i >>> 5] &= ~(1 << (i & 31))
      if (this._seg[h][i >>> 5] > 0) {
        break
      }
      i >>>= 5
    }
  }

  /**
   * 返回x右侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回null.
   */
  next(i: number): number | null {
    if (i < 0) {
      i = 0
    }
    if (i >= this._n) {
      return null
    }

    for (let h = 0; h < this._lg; h++) {
      if (i >>> 5 === this._seg[h].length) {
        break
      }
      let d = this._seg[h][i >>> 5] >>> (i & 31)
      if (!d) {
        i = (i >>> 5) + 1
        continue
      }
      // !trailingZeros32: 31 - Math.clz32(x & -x)
      i += 31 - Math.clz32(d & -d)
      for (let g = h - 1; ~g; g--) {
        i <<= 5
        const tmp = this._seg[g][i >>> 5]
        i += 31 - Math.clz32(tmp & -tmp)
      }
      return i
    }

    return null
  }

  /**
   * 返回x左侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回null.
   */
  prev(i: number): number | null {
    if (i < 0) {
      return null
    }
    if (i >= this._n) {
      i = this._n - 1
    }

    for (let h = 0; h < this._lg; h++) {
      if (i === -1) {
        break
      }
      let d = this._seg[h][i >>> 5] << (31 - (i & 31))
      if (!d) {
        i = (i >>> 5) - 1
        continue
      }

      i -= Math.clz32(d)
      for (let g = h - 1; ~g; g--) {
        i <<= 5
        i += 31 - Math.clz32(this._seg[g][i >>> 5])
      }
      return i
    }

    return null
  }

  /**
   * 遍历[start,end)区间内的元素.
   */
  enumerateRange(start: number, end: number, f: (v: number) => void): void {
    let x: number | null = start - 1
    while (true) {
      x = this.next(x + 1)
      if (x == null || x >= end) {
        break
      }
      f(x)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateRange(0, this._n, v => sb.push(v.toString()))
    return `Finder(${sb.join(', ')})`
  }
}

export { Finder }
