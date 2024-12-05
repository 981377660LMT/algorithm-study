/* eslint-disable implicit-arrow-linebreak */
/* eslint-disable no-param-reassign */
/* eslint-disable no-inner-declarations */
/* eslint-disable class-methods-use-this */

/**
 * @summary
 * 高效计算数列的前缀和，区间和
 * 树状数组或二叉索引树（Binary Indexed Tree, Fenwick Tree）
 * 性质
 * 1. tree[x]保存以x为根的子树中叶节点值的和
 * 2. tree[x]的父节点为tree[x+lowbit(x)]
 * 3. tree[x]节点覆盖的长度等于lowbit(x)
 * 4. 树的高度为logn+1
 */

// 下标从0开始
// 1.BITArray: 单点修改, 区间查询
// 2.BITMap: 单点修改, 区间查询
// 3.BITRangeAddPointGetArray: 区间修改, 单点查询(差分)
// 4.BITRangeAddPointGetMap: 区间修改, 单点查询(差分)
// 5.BITRangeAddRangeSumArray: 区间修改, 区间查询
// 6.BITRangeAddRangeSumMap: 区间修改, 区间查询
// 7.BITPrefixArray: 单点修改, 前缀查询
// 8.BITPrefixMap: 单点修改, 前缀查询

class BITArray {
  readonly n: number
  private readonly _data: number[]
  private _total = 0

  constructor(n: number, f?: (i: number) => number) {
    if (f == undefined) {
      this.n = n
      this._data = Array(n).fill(0)
    } else {
      this.n = n
      this._data = Array(n)
      for (let i = 0; i < n; i++) {
        this._data[i] = f(i)
        this._total += this._data[i]
      }
      for (let i = 1; i <= n; i++) {
        let j = i + (i & -i)
        if (j <= n) {
          this._data[j - 1] += this._data[i - 1]
        }
      }
    }
  }

  add(index: number, v: number): void {
    this._total += v
    index += 1
    while (index <= this.n) {
      this._data[index - 1] += v
      index += index & -index
    }
  }

  queryPrefix(end: number): number {
    if (end > this.n) {
      end = this.n
    }
    let res = 0
    while (end > 0) {
      res += this._data[end - 1]
      end -= end & -end
    }
    return res
  }

  queryRange(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this.n) end = this.n
    if (start >= end) return 0
    if (start === 0) return this.queryPrefix(end)
    let pos = 0
    let neg = 0
    while (end > start) {
      pos += this._data[end - 1]
      end &= end - 1
    }
    while (start > end) {
      neg += this._data[start - 1]
      start &= start - 1
    }
    return pos - neg
  }

  queryAll(): number {
    return this._total
  }

  /**
   * 查询满足`predicate`的最大的`end`(不包含).
   */
  maxRight(predicate: (index: number, preSum: number) => boolean): number {
    let i = 0
    let s = 0
    let k = 1
    while (2 * k <= this.n) {
      k *= 2
    }
    while (k > 0) {
      if (i + k - 1 < this.n) {
        let t = s + this._data[i + k - 1]
        if (predicate(i + k, t)) {
          i += k
          s = t
        }
      }
      k >>>= 1
    }
    return i
  }

  /**
   * 01树状数组查找第 k(0-based) 个1的位置.
   */
  kth(k: number): number {
    return this.maxRight((_, preSum) => preSum <= k)
  }

  toString(): string {
    return `BITArray: [${Array.from({ length: this.n }, (_, i) => this.queryRange(i, i + 1)).join(', ')}]`
  }
}

class BITMap {
  private readonly _n: number
  private readonly _data: Map<number, number> = new Map()
  private _total = 0

  constructor(n: number) {
    if (n > 2 ** 31 - 1) throw new Error('BITMap: n must be less than 2^31-1')
    this._n = n
  }

  add(index: number, v: number): void {
    this._total += v
    index += 1
    while (index <= this._n) {
      this._data.set(index, (this._data.get(index) || 0) + v)
      index += index & -index
    }
  }

  queryPrefix(end: number): number {
    if (end > this._n) {
      end = this._n
    }
    let res = 0
    while (end > 0) {
      res += this._data.get(end) || 0
      end -= end & -end
    }
    return res
  }

  queryRange(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return 0
    if (start === 0) return this.queryPrefix(end)
    let pos = 0
    let neg = 0
    while (end > start) {
      pos += this._data.get(end) || 0
      end &= end - 1
    }
    while (start > end) {
      neg += this._data.get(start) || 0
      start &= start - 1
    }
    return pos - neg
  }

  queryAll(): number {
    return this._total
  }

  /**
   * 查询满足`predicate`的最大的`end`(不包含).
   */
  maxRight(predicate: (index: number, preSum: number) => boolean): number {
    let i = 0
    let s = 0
    let k = 1
    while (2 * k <= this._n) {
      k *= 2
    }
    while (k > 0) {
      if (i + k - 1 < this._n) {
        let t = s + (this._data.get(i + k) || 0)
        if (predicate(i + k, t)) {
          i += k
          s = t
        }
      }
      k >>>= 1
    }
    return i
  }

  /**
   * 01树状数组查找第 k(0-based) 个1的位置.
   */
  kth(k: number): number {
    return this.maxRight((_, preSum) => preSum <= k)
  }
}

class BITRangeAddPointGetArray {
  private readonly _bit: BITArray

  constructor(n: number, f?: (i: number) => number) {
    this._bit = new BITArray(n, f)
  }

  addRange(start: number, end: number, delta: number): void {
    this._bit.add(start, delta)
    this._bit.add(end, -delta)
  }

  get(index: number): number {
    return this._bit.queryPrefix(index + 1)
  }

  toString(): string {
    return `BITRangeAddPointGetArray: [${Array.from({ length: this._bit.n }, (_, i) => this.get(i)).join(', ')}]`
  }
}

class BITRangeAddRangeSumArray {
  private readonly _n: number
  private readonly _bit0: BITArray
  private readonly _bit1: BITArray

  constructor(n: number, f?: (i: number) => number) {
    if (f == undefined) {
      this._n = n
      this._bit0 = new BITArray(n)
      this._bit1 = new BITArray(n)
    } else {
      this._n = n
      this._bit0 = new BITArray(n, f)
      this._bit1 = new BITArray(n)
    }
  }

  add(index: number, delta: number): void {
    this._bit0.add(index, delta)
  }

  addRange(start: number, end: number, delta: number): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    this._bit0.add(start, -delta * start)
    this._bit0.add(end, delta * end)
    this._bit1.add(start, delta)
    this._bit1.add(end, -delta)
  }

  queryRange(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return 0
    let rightRes = this._bit1.queryPrefix(end) * end + this._bit0.queryPrefix(end)
    let leftRes = this._bit1.queryPrefix(start) * start + this._bit0.queryPrefix(start)
    return rightRes - leftRes
  }

  toString(): string {
    return `BITRangeAddRangeSumArray: [${Array.from({ length: this._n }, (_, i) => this.queryRange(i, i + 1)).join(', ')}]`
  }
}

class BITRangeAddPointGetMap {
  private readonly _bit: BITMap

  constructor(n: number) {
    this._bit = new BITMap(n)
  }

  addRange(start: number, end: number, delta: number): void {
    this._bit.add(start, delta)
    this._bit.add(end, -delta)
  }

  get(index: number): number {
    return this._bit.queryPrefix(index + 1)
  }
}

class BITRangeAddRangeSumMap {
  private readonly _n: number
  private readonly _bit0: BITMap
  private readonly _bit1: BITMap

  constructor(n: number) {
    this._n = n
    this._bit0 = new BITMap(n)
    this._bit1 = new BITMap(n)
  }

  add(index: number, delta: number): void {
    this._bit0.add(index, delta)
  }

  addRange(start: number, end: number, delta: number): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    this._bit0.add(start, -delta * start)
    this._bit0.add(end, delta * end)
    this._bit1.add(start, delta)
    this._bit1.add(end, -delta)
  }

  queryRange(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return 0
    let rightRes = this._bit1.queryPrefix(end) * end + this._bit0.queryPrefix(end)
    let leftRes = this._bit1.queryPrefix(start) * start + this._bit0.queryPrefix(start)
    return rightRes - leftRes
  }

  toString(): string {
    return `BITRangeAddRangeSumMap: [${Array.from({ length: this._n }, (_, i) => this.queryRange(i, i + 1)).join(', ')}]`
  }
}

class BITPrefixArray<S> {
  private readonly _n: number
  private readonly _data: S[]
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S

  constructor(n: number, e: () => S, op: (a: S, b: S) => S, f?: (i: number) => S) {
    this._n = n
    this._data = Array(n)
    for (let i = 0; i < n; i++) this._data[i] = e()
    this._e = e
    this._op = op
    if (f != undefined) {
      for (let i = 1; i <= n; i++) {
        let j = i + (i & -i)
        if (j <= n) {
          this._data[j - 1] = op(this._data[j - 1], this._data[i - 1])
        }
      }
    }
  }

  update(index: number, value: S): void {
    index += 1
    while (index <= this._n) {
      this._data[index - 1] = this._op(this._data[index - 1], value)
      index += index & -index
    }
  }

  query(end: number): S {
    if (end > this._n) {
      end = this._n
    }
    let res = this._e()
    while (end > 0) {
      res = this._op(res, this._data[end - 1])
      end -= end & -end
    }
    return res
  }
}

class BITPrefixMap<S> {
  private readonly _n: number
  private readonly _data: Map<number, S> = new Map()
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S

  constructor(n: number, e: () => S, op: (a: S, b: S) => S) {
    if (n > 2 ** 31 - 1) throw new Error('BITPrefixMap: n must be less than 2^31-1')
    this._n = n
    this._e = e
    this._op = op
  }

  update(index: number, value: S): void {
    index += 1
    while (index <= this._n) {
      this._data.set(index - 1, this._op(this._data.get(index - 1) ?? this._e(), value))
      index += index & -index
    }
  }

  query(end: number): S {
    if (end > this._n) {
      end = this._n
    }
    let res = this._e()
    while (end > 0) {
      res = this._op(res, this._data.get(end - 1) ?? this._e())
      end -= end & -end
    }
    return res
  }
}

export {
  BITArray,
  BITMap,
  BITRangeAddPointGetArray,
  BITRangeAddPointGetMap,
  BITRangeAddRangeSumArray,
  BITRangeAddRangeSumMap,
  BITPrefixArray,
  BITPrefixMap
}

if (require.main === module) {
  const bitArray = new BITArray(10)
  console.log(bitArray.toString())
  const bitMap = new BITMap(1e9 + 10)
  bitMap.add(1e7, 11)
  console.log(bitMap.queryPrefix(1e7 + 1))

  const bitRangeAddRangeSumArray = new BITRangeAddRangeSumArray(10)
  bitRangeAddRangeSumArray.addRange(1, 3, 2)
  console.log(bitRangeAddRangeSumArray.toString())
  const bitRangeAddRangeSumMap = new BITRangeAddRangeSumMap(10)
  bitRangeAddRangeSumMap.addRange(1, 3, 2)
  console.log(bitRangeAddRangeSumMap.queryRange(1, 3))

  const bitPrefixArray = new BITPrefixArray(10, () => 0, Math.max)
  bitPrefixArray.update(0, 1)
  bitPrefixArray.update(1, 2)
  console.log(bitPrefixArray.query(2))
  const bitPrefixMap = new BITPrefixMap(1e9, () => 0, Math.max)
  bitPrefixMap.update(0, 1)
  bitPrefixMap.update(1, 2)
  console.log(bitPrefixMap.query(2))

  const bb = new BITRangeAddPointGetArray(10)
  bb.addRange(1, 3, 2)
  console.log(bb.toString(), bb.get(1))

  // https://leetcode.cn/problems/maximum-white-tiles-covered-by-a-carpet/
  function maximumWhiteTiles(tiles: number[][], carpetLen: number): number {
    const bit = new BITRangeAddRangeSumMap(1e9 + 10)
    let res = 0
    tiles.forEach(([left, right]) => {
      bit.addRange(left, right + 1, 1)
    })
    tiles.forEach(([left]) => {
      res = Math.max(res, bit.queryRange(left, left + carpetLen))
    })
    return res
  }
}
