/* eslint-disable no-param-reassign */
/* eslint-disable no-constant-condition */

/**
 * 利用位运算寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 */
class FastSet {
  private readonly _n: number
  private readonly _lg: number
  private readonly _seg: Uint32Array[]
  private _size = 0

  /**
   * @param n [0, n).
   * @param initValue 初始化时是否全部置为1.默认初始时所有位置都未被访问过.
   */
  constructor(n: number, initValue = false) {
    this._n = n
    const seg: Uint32Array[] = []
    while (true) {
      if (initValue) {
        seg.push(new Uint32Array((n + 31) >>> 5).fill(-1))
      } else {
        seg.push(new Uint32Array((n + 31) >>> 5))
      }
      n = (n + 31) >>> 5
      if (n <= 1) {
        break
      }
    }
    this._lg = seg.length
    this._seg = seg
  }

  insert(i: number): boolean {
    if (this.has(i)) return false
    for (let h = 0; h < this._lg; h++) {
      this._seg[h][i >>> 5] |= 1 << (i & 31)
      i >>>= 5
    }
    this._size++
    return true
  }

  has(i: number): boolean {
    return !!(this._seg[0][i >>> 5] & (1 << (i & 31)))
  }

  erase(i: number): boolean {
    if (!this.has(i)) return false
    for (let h = 0; h < this._lg; h++) {
      const cache = this._seg[h]
      cache[i >>> 5] &= ~(1 << (i & 31))
      if (cache[i >>> 5]) break
      i >>>= 5
    }
    this._size--
    return true
  }

  /**
   * 返回x右侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回`n`.
   */
  next(i: number): number {
    if (i < 0) i = 0
    if (i >= this._n) return this._n

    for (let h = 0; h < this._lg; h++) {
      const cacheH = this._seg[h]
      if (i >>> 5 === cacheH.length) break
      let d = cacheH[i >>> 5] >>> (i & 31)
      if (d === 0) {
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

    return this._n
  }

  /**
   * 返回x左侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回`-1`.
   */
  prev(i: number): number {
    if (i < 0) return -1
    if (i >= this._n) i = this._n - 1

    for (let h = 0; h < this._lg; h++) {
      if (i === -1) break
      let d = this._seg[h][i >>> 5] << (31 - (i & 31))
      if (d === 0) {
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

    return -1
  }

  /**
   * 遍历[start,end)区间内的元素.
   */
  enumerateRange(start: number, end: number, f: (v: number) => void): void {
    let x = start - 1
    while (true) {
      x = this.next(x + 1)
      if (x >= end) break
      f(x)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateRange(0, this._n, v => sb.push(v.toString()))
    return `FastSet(${this.size}){${sb.join(', ')}}`
  }

  get min(): number {
    return this.next(-1)
  }

  get max(): number {
    return this.prev(this._n)
  }

  get size(): number {
    return this._size
  }
}

if (require.main === module) {
  const s = new FastSet(100)
  s.insert(1)
  console.log(s.toString(), s.prev(1), s.next(1), s.prev(0), s.next(0))

  const n = 1e7
  const set2 = new FastSet(n)
  console.time('FastSet')
  for (let i = 0; i < n; i++) {
    set2.insert(i)
    set2.next(i)
    set2.prev(i)
    set2.has(i)
    set2.erase(i)
    set2.insert(i)
  }
  console.timeEnd('FastSet') // 270ms

  const allOne = new FastSet(10, true)
  console.log(allOne.toString())
}

export { FastSet }
