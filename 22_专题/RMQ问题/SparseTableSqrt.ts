// https://www.cnblogs.com/MoyouSayuki/p/17595714.html
// https://www.luogu.com.cn/problem/solution/P3793
// https://kewth.github.io/2019/10/11/RMQ/
//
// 块间分为整块和散块，对于散块可以预处理出每一个块的前后缀最大值，
// 这样预处理是 (O(n)) 的，查询降为 (O(1))，
// 对于整块可以整一个 Sparse Table，把每一个块的最大值看成一个元素，维护块的最大值的最大值，
// 这样可以做到预处理 (O(sqrt nlog(sqrt n)))，询问 (O(1))。
// !分块优化ST表，大概就是把ST表分块，然后统计每一块的前后缀最大值，
// 就可以在O(1∼ sqrt(n))的时间里完成查询并做到`节省空间`的效果，
// 这种方法的应用空间很广泛，甚至可以拓展到所有有结合律的函数。
//
// 瓶颈：左右端点恰好在同一个块中，此时只能遍历块求解

import { SparseTable } from './SparseTable'

class SparseTableSqrt<S> {
  private readonly _arr: ArrayLike<S>
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S
  private readonly _belong: (index: number) => number
  private readonly _st: SparseTable<S>
  private readonly _pre: S[]
  private readonly _suf: S[]

  constructor(
    arr: ArrayLike<S>,
    e: () => S,
    op: (a: S, b: S) => S,
    blockSize = (Math.sqrt(arr.length) + 1) | 0
  ) {
    const n = arr.length
    const belong = (index: number) => (index / blockSize) | 0
    const blockStart = (index: number) => index * blockSize
    const blockEnd = (index: number) => Math.min((index + 1) * blockSize, n)
    const blockCount = 1 + ((n / blockSize) | 0)

    const blockRes: S[] = Array(blockCount)
    for (let i = 0; i < blockCount; i++) blockRes[i] = e()
    for (let i = 0; i < n; i++) {
      const bid = belong(i)
      blockRes[bid] = op(blockRes[bid], arr[i])
    }
    const st = new SparseTable(blockRes, e, op)

    const pre: S[] = Array(n)
    for (let bid = 0; bid < blockCount; bid++) {
      let res = e()
      for (let i = blockStart(bid); i < blockEnd(bid); i++) {
        res = op(res, arr[i])
        pre[i] = res
      }
    }

    const suf: S[] = Array(n)
    for (let i = 0; i < n; i++) suf[i] = e()
    for (let bid = 0; bid < blockCount; bid++) {
      let res = e()
      for (let i = blockEnd(bid) - 1; i >= blockStart(bid); i--) {
        res = op(arr[i], res)
        suf[i] = res
      }
    }

    this._arr = arr
    this._e = e
    this._op = op
    this._belong = belong
    this._st = st
    this._pre = pre
    this._suf = suf
  }

  /**
   * 查询左闭右开区间`[start, end)`的贡献值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): S {
    if (start < 0) start = 0
    if (end > this._arr.length) end = this._arr.length
    if (start >= end) return this._e()

    const bid1 = this._belong(start)
    const bid2 = this._belong(end - 1)
    if (bid1 === bid2) {
      let res = this._e()
      for (let i = start; i < end; i++) res = this._op(res, this._arr[i])
      return res
    }

    let res = this._suf[start]
    res = this._op(res, this._st.query(bid1 + 1, bid2))
    res = this._op(res, this._pre[end - 1])
    return res
  }
}

class SparseTableSqrtInt32 {
  private readonly _arr: ArrayLike<number>
  private readonly _e: () => number
  private readonly _op: (a: number, b: number) => number
  private readonly _belong: (index: number) => number
  private readonly _st: SparseTable // 数组较短时使用 `SparseTable` 更好
  private readonly _pre: Int32Array
  private readonly _suf: Int32Array

  constructor(
    arr: ArrayLike<number>,
    e: () => number,
    op: (a: number, b: number) => number,
    blockSize = (Math.sqrt(arr.length) + 1) | 0
  ) {
    const n = arr.length
    const belong = (index: number) => (index / blockSize) | 0
    const blockStart = (index: number) => index * blockSize
    const blockEnd = (index: number) => Math.min((index + 1) * blockSize, n)
    const blockCount = 1 + ((n / blockSize) | 0)

    const blockRes = Array(blockCount).fill(e())
    for (let i = 0; i < n; i++) {
      const bid = belong(i)
      blockRes[bid] = op(blockRes[bid], arr[i])
    }
    const st = new SparseTable(blockRes, e, op)

    const pre = new Int32Array(n)
    for (let bid = 0; bid < blockCount; bid++) {
      let res = e()
      for (let i = blockStart(bid); i < blockEnd(bid); i++) {
        res = op(res, arr[i])
        pre[i] = res
      }
    }

    const suf = new Int32Array(n)
    for (let i = 0; i < n; i++) suf[i] = e()
    for (let bid = 0; bid < blockCount; bid++) {
      let res = e()
      for (let i = blockEnd(bid) - 1; i >= blockStart(bid); i--) {
        res = op(arr[i], res)
        suf[i] = res
      }
    }

    this._arr = arr
    this._e = e
    this._op = op
    this._belong = belong
    this._st = st
    this._pre = pre
    this._suf = suf
  }

  /**
   * 查询左闭右开区间`[start, end)`的贡献值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._arr.length) end = this._arr.length
    if (start >= end) return this._e()

    const bid1 = this._belong(start)
    const bid2 = this._belong(end - 1)
    if (bid1 === bid2) {
      let res = this._e()
      for (let i = start; i < end; i++) res = this._op(res, this._arr[i])
      return res
    }

    let res = this._suf[start]
    res = this._op(res, this._st.query(bid1 + 1, bid2))
    res = this._op(res, this._pre[end - 1])
    return res
  }
}

export { SparseTableSqrt, SparseTableSqrtInt32 }

if (require.main === module) {
  const st = new SparseTableSqrt([1, 2, 3, 4, 5, 6, 7, 8, 9], () => 0, Math.max)
  console.log(st.query(0, 5))
  // console.log(st.query(0, 4))

  const naive = (arr: number[], start: number, end: number) => {
    let res = 0
    for (let i = start; i < end; i++) {
      res = Math.max(res, arr[i])
    }
    return res
  }

  const test = (arr: number[], start: number, end: number) => {
    const st = new SparseTableSqrt(arr, () => 0, Math.max)
    const res = st.query(start, end)
    const ans = naive(arr, start, end)
    if (res !== ans) {
      throw new Error()
    }
  }

  const arr = [1, 2, 3, 4, 5, 6, 7, 8, 9]
  for (let i = 0; i < arr.length; i++) {
    for (let j = i + 1; j <= arr.length; j++) {
      test(arr, i, j)
    }
  }
}
