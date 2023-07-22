/* eslint-disable no-param-reassign */

/**
 * 支持幺半群的静态区间查询.
 * 比`st表`稍慢一点,但是可以不用满足幂等性.
 * @deprecated 使用线段树代替.
 */
class DisjointSparseTable<S> {
  private readonly _n: number
  private readonly _data: S[]
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S

  constructor(leaves: ArrayLike<S>, e: () => S, op: (a: S, b: S) => S) {
    const n = leaves.length
    let log = 1
    while (1 << log < n) log++
    const data: S[] = Array(log * n)
    for (let i = 0; i < n; i++) data[i] = leaves[i]
    for (let i = 1; i < log; i++) {
      for (let j = 0; j < n; j++) data[i * n + j] = leaves[j]
      const b = 1 << i
      for (let m = b; m <= n; m += 2 * b) {
        const l = m - b
        const r = Math.min(m + b, n)
        for (let j = m - 1; j >= l + 1; j--) {
          const pos = i * n + j
          data[pos - 1] = op(data[pos - 1], data[pos])
        }
        for (let j = m; j < r - 1; j++) {
          const pos = i * n + j
          data[pos + 1] = op(data[pos], data[pos + 1])
        }
      }
    }

    this._n = n
    this._data = data
    this._e = e
    this._op = op
  }

  /**
   * [start,end).
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): S {
    if (start >= end) return this._e()
    end--
    if (start === end) return this._data[start]
    const k = 31 - Math.clz32(start ^ end)
    return this._op(this._data[k * this._n + start], this._data[k * this._n + end])
  }

  /**
   * 返回最大的`right`使得`[left,right)`内的值满足`check`.
   * 0 <= left <= right <= n.
   */
  maxRight(left: number, check: (s: S) => boolean): number {
    if (left === this._n) return this._n
    let ok = left
    let ng = this._n + 1
    while (ok + 1 < ng) {
      const mid = (ok + ng) >> 1
      if (check(this.query(left, mid))) {
        ok = mid
      } else {
        ng = mid
      }
    }
    return ok
  }

  /**
   * 返回最小的`left`使得`[left,right)`内的值满足`check`.
   * 0 <= left <= right <= n.
   */
  minLeft(right: number, check: (s: S) => boolean): number {
    if (!right) return 0
    let ok = right
    let ng = -1
    while (ng + 1 < ok) {
      const mid = (ok + ng) >> 1
      if (check(this.query(mid, right))) {
        ok = mid
      } else {
        ng = mid
      }
    }
    return ok
  }
}

export { DisjointSparseTable }

if (require.main === module) {
  const nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  const st = new DisjointSparseTable(
    nums,
    () => 0,
    (a, b) => a + b
  )
  console.log(
    st.query(0, 10),
    st.maxRight(0, e => e < 10),
    st.minLeft(10, e => e < 20)
  )
}
