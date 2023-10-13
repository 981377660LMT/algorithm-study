/* eslint-disable prefer-destructuring */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */
/* eslint-disable no-inner-declarations */

// 珂朵莉树(ODT)/Intervals
// !noneValue不使用symbol,而是自定义的哨兵值,更加灵活.

/**
 * 珂朵莉树，基于数据随机的颜色段均摊。
 * `FastSet`实现.
 */
class ODT<S> {
  private _len = 0
  private _count = 0
  private readonly _leftLimit: number
  private readonly _rightLimit: number
  private readonly _noneValue: S
  private readonly _data: S[]
  private readonly _fs: FastSet

  /**
   * 指定区间长度和哨兵值建立一个ODT.初始时,所有位置的值为 {@link noneValue}.
   * @param n 区间范围为`[0, n)`.
   * @param noneValue 表示空值的哨兵值.
   */
  constructor(n: number, noneValue: S) {
    const data = Array(n)
    for (let i = 0; i < n; i++) data[i] = noneValue
    const fs = new FastSet(n)
    fs.insert(0)

    this._leftLimit = 0
    this._rightLimit = n
    this._noneValue = noneValue
    this._data = data
    this._fs = fs
  }

  /**
   * 返回包含`x`的区间的信息.
   * 0 <= x < n.
   */
  get(x: number, erase = false): [start: number, end: number, value: S] | undefined {
    if (x < this._leftLimit || x >= this._rightLimit) return undefined
    const start = this._fs.prev(x)
    const end = this._fs.next(x + 1)
    const value = this._data[start]
    if (erase && value !== this._noneValue) {
      this._len--
      this._count -= end - start
      this._data[start] = this._noneValue
      this._mergeAt(start)
      this._mergeAt(end)
    }
    return [start, end, value]
  }

  set(start: number, end: number, value: S): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerateRange(start, end, () => {}, true) // remove
    this._fs.insert(start)
    this._data[start] = value
    if (value !== this._noneValue) {
      this._len++
      this._count += end - start
    }
    this._mergeAt(start)
    this._mergeAt(end)
  }

  enumerateAll(f: (start: number, end: number, value: S) => void): void {
    this.enumerateRange(this._leftLimit, this._rightLimit, f, false)
  }

  /**
   * 遍历范围`[start, end)`内的所有区间.
   */
  enumerateRange(
    start: number,
    end: number,
    f: (start: number, end: number, value: S) => void,
    erase = false
  ): void {
    if (start < this._leftLimit) start = this._leftLimit
    if (end > this._rightLimit) end = this._rightLimit
    if (start >= end) return

    const none = this._noneValue
    if (!erase) {
      let left = this._fs.prev(start)
      while (left < end) {
        const right = this._fs.next(left + 1)
        f(Math.max(left, start), Math.min(right, end), this._data[left])
        left = right
      }
      return
    }

    let p = this._fs.prev(start)
    if (p < start) {
      this._fs.insert(start)
      const v = this._data[p]
      this._data[start] = v
      if (v !== none) {
        this._len++
      }
    }

    p = this._fs.next(end)
    if (end < p) {
      const v = this._data[this._fs.prev(end)]
      this._data[end] = v
      this._fs.insert(end)
      if (v !== none) {
        this._len++
      }
    }

    p = start
    while (p < end) {
      const q = this._fs.next(p + 1)
      const x = this._data[p]
      f(p, q, x)
      if (x !== none) {
        this._len--
        this._count -= q - p
      }
      this._fs.erase(p)
      p = q
    }

    this._fs.insert(start)
    this._data[start] = none
  }

  toString(): string {
    const sb: string[] = [`ODT(${this.length}) {`]
    this.enumerateAll((start, end, value) => {
      const v = value === this._noneValue ? 'null' : value
      sb.push(`  [${start},${end}):${v}`)
    })
    sb.push('}')
    return sb.join('\n')
  }

  /**
   * 区间个数.
   */
  get length(): number {
    return this._len
  }

  /**
   * 区间内元素个数之和.
   */
  get count(): number {
    return this._count
  }

  private _mergeAt(p: number): void {
    if (p <= 0 || this._rightLimit <= p) return
    const q = this._fs.prev(p - 1)
    const dataP = this._data[p]
    const dataQ = this._data[q]
    if (dataP === dataQ) {
      if (dataP !== this._noneValue) this._len--
      this._fs.erase(p)
    }
  }
}

/**
 * 利用位运算寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 * 初始时,所有位置都未被访问过.
 */
class FastSet {
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
      if (this._seg[h][i >>> 5]) {
        break
      }
      i >>>= 5
    }
  }

  /**
   * 返回x右侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回`n`.
   */
  next(i: number): number {
    if (i < 0) {
      i = 0
    }
    if (i >= this._n) {
      return this._n
    }

    for (let h = 0; h < this._lg; h++) {
      if (i >>> 5 === this._seg[h].length) {
        break
      }
      let d = this._seg[h][i >>> 5] >>> (i & 31)
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
    if (i < 0) {
      return -1
    }
    if (i >= this._n) {
      i = this._n - 1
    }

    for (let h = 0; h < this._lg; h++) {
      if (i === -1) {
        break
      }
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
      if (x >= end) {
        break
      }
      f(x)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateRange(0, this._n, v => sb.push(v.toString()))
    return `FastSet(${sb.join(', ')})`
  }

  get min(): number | null {
    return this.next(-1)
  }

  get max(): number | null {
    return this.prev(this._n)
  }
}

export { ODT }

if (require.main === module) {
  // const INF = 2e15
  // const odt = new ODT(10, INF)
  // console.log(odt.toString())
  // odt.set(0, 10, 1)
  // odt.set(2, 5, 2)
  // console.log(odt.get(8))
  // console.log(odt.toString())
  // odt.enumerateRange(
  //   1,
  //   7,
  //   (start, end, value) => {
  //     console.log(start, end, value)
  //   },
  //   true
  // )
  // console.log(odt.toString(), odt.length)

  // 352. 将数据流变为多个不相交区间
  // https://leetcode.cn/problems/data-stream-as-disjoint-intervals/
  class SummaryRanges {
    private readonly _odt = new ODT(1e4 + 10, -1)

    addNum(value: number): void {
      this._odt.set(value, value + 1, 0)
    }

    getIntervals(): number[][] {
      const res: number[][] = []
      this._odt.enumerateAll((start, end, value) => {
        if (value === 0) res.push([start, end - 1])
      })
      return res
    }
  }
}
