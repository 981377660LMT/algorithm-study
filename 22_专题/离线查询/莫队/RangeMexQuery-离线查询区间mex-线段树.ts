const INF = 2e15

/**
 * 离线查询区间mex.基于线段树实现.
 */
class RangeMexQuery {
  private readonly _nums: ArrayLike<number>
  private readonly _query: [start: number, end: number][] = []

  constructor(nums: ArrayLike<number>) {
    this._nums = nums
  }

  /**
   * [start, end).
   * 0 <= start <= end <= n.
   */
  addQuery(start: number, end: number): void {
    this._query.push([start, end])
  }

  /**
   * @param mexStart mex的起始值(从0开始还是从1开始).
   */
  run(mexStart: number): number[] {
    const n = this._nums.length
    const leaves = Array(n + 2).fill(-1)
    const seg = new _SegTree(leaves)
    const q = this._query.length
    const res = Array(q)
    const ids: number[][] = Array(n + 1)
    for (let i = 0; i < ids.length; i++) ids[i] = []
    this._query.forEach(([_, end], i) => {
      ids[end].push(i)
    })

    for (let i = 0; i < n + 1; i++) {
      ids[i].forEach(q => {
        const start = this._query[q][0]
        const mex = seg.maxRight(mexStart, x => x >= start)
        res[q] = mex
      })
      if (i < n && this._nums[i] < n + 2) {
        seg.set(this._nums[i], i)
      }
    }

    return res
  }
}

class _SegTree {
  private readonly _n: number
  private readonly _size: number
  private readonly _seg: number[]

  constructor(leaves: number[]) {
    const n = leaves.length
    let size = 1
    while (size < n) size <<= 1
    const seg = Array(size << 1)
    for (let i = 0; i < n; i++) seg[i + size] = leaves[i]
    for (let i = size - 1; ~i; i--) seg[i] = Math.min(seg[i << 1], seg[(i << 1) | 1])
    this._n = n
    this._size = size
    this._seg = seg
  }

  get(index: number): number {
    if (index < 0 || index >= this._n) return INF
    return this._seg[index + this._size]
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    this._seg[index] = value
    for (index >>= 1; index > 0; index >>= 1) {
      this._seg[index] = Math.min(this._seg[index << 1], this._seg[(index << 1) | 1])
    }
  }

  maxRight(left: number, predicate: (value: number) => boolean): number {
    if (left === this._n) return this._n
    left += this._size
    let res = INF
    while (true) {
      if (!(left & 1)) left >>= 1
      if (!predicate(Math.min(res, this._seg[left]))) {
        while (left < this._size) {
          left <<= 1
          const tmp = Math.min(res, this._seg[left])
          if (predicate(tmp)) {
            res = tmp
            left++
          }
        }
        return left - this._size
      }
      res = Math.min(res, this._seg[left])
      left++
      if ((left & -left) === left) break
    }
    return this._n
  }
}

export { RangeMexQuery }

if (require.main === module) {
  const nums = [1, 2, 3, 5, 6, 6, 7, 8, 9]
  const M = new RangeMexQuery(nums)
  M.addQuery(0, 3)
  M.addQuery(0, 4)
  console.log(M.run(1))

  const n = 1e5
  const arr = Array(n)
    .fill(0)
    .map((_, i) => i + 1)
  const M2 = new RangeMexQuery(arr)
  for (let i = 0; i < n; i++) {
    M2.addQuery(0, i)
  }
  console.time('mex')
  M2.run(1)
  console.timeEnd('mex')
}
