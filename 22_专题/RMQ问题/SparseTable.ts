/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

import assert from 'assert'

const INF = 2e15

/**
 * st表适用于可重复贡献问题/RMQ静态区间最值查询。
 * 可重复贡献问题是指对于运算`opt`，满足`x opt x = x`，则对应的区间询问就是一个可重复贡献问题。
 * !幂等性
 * 例如，最大值有 `max(a,a) = a`，gcd有 `gcd(a,a) = a` ，所以RMQ和区间GCD就是一个可重复贡献问题。
 * 像区间和就不具有这个性质，如果求区间和的时候采用的预处理区间重叠了，
 * 则会导致重曼部分被计算两次，这是我们所不愿意看到的。
 * 另外，opt还必须满足结合律才能使用ST表求解。
 * @see {@link https://oi-wiki.org/ds/sparse-table/}
 * @important
 * - cpu-cache optimized `dp[bit][n]`
 * - compress to 1d `dp[bit][n] -> dp[bit*n]`
 */
class SparseTable<S = number> {
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S
  private readonly _dp: S[]
  private readonly _n: number

  constructor(nums: ArrayLike<S>, e: () => S, op: (a: S, b: S) => S) {
    const n = nums.length
    const size = 32 - Math.clz32(n)
    this._n = n
    this._e = e
    this._op = op
    this._dp = Array(size * n) // !dp[i][j]表示区间[j,j+2**i)的贡献值
    for (let i = 0; i < n; i++) this._dp[i] = nums[i]
    for (let i = 1; i < size; i++) {
      for (let j = 0; j < n - (1 << i) + 1; j++) {
        const pos = i * n + j
        this._dp[pos] = op(this._dp[pos - n], this._dp[pos - n + (1 << (i - 1))])
      }
    }
  }

  /**
   * 查询左闭右开区间`[start, end)`的贡献值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): S {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return this._e()
    const k = 31 - Math.clz32(end - start)
    return this._op(this._dp[k * this._n + start], this._dp[k * this._n + end - (1 << k)])
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

/**
 * 元素类型为`int32`类型的st表.
 */
class SparseTableInt32 {
  private readonly _e: () => number
  private readonly _op: (a: number, b: number) => number
  private readonly _dp: Int32Array
  private readonly _n: number

  constructor(nums: ArrayLike<number>, e: () => number, op: (a: number, b: number) => number) {
    const n = nums.length
    const size = 32 - Math.clz32(n)
    this._n = n
    this._e = e
    this._op = op
    this._dp = new Int32Array(size * n) // !dp[i][j]表示区间[j,j+2**i)的贡献值
    for (let i = 0; i < n; i++) this._dp[i] = nums[i]
    for (let i = 1; i < size; i++) {
      for (let j = 0; j < n - (1 << i) + 1; j++) {
        const pos = i * n + j
        this._dp[pos] = op(this._dp[pos - n], this._dp[pos - n + (1 << (i - 1))])
      }
    }
  }

  /**
   * 查询左闭右开区间`[start, end)`的贡献值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return this._e()
    const k = 31 - Math.clz32(end - start)
    return this._op(this._dp[k * this._n + start], this._dp[k * this._n + end - (1 << k)])
  }

  /**
   * 返回最大的`right`使得`[left,right)`内的值满足`check`.
   * 0 <= left <= right <= n.
   */
  maxRight(left: number, check: (s: number) => boolean): number {
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
  minLeft(right: number, check: (s: number) => boolean): number {
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

export { SparseTable, SparseTableInt32 }

if (require.main === module) {
  const st = new SparseTable(
    [9, 12, 3, 7, 15],
    () => 0,
    (a, b) => Math.max(a, b)
  )
  assert.strictEqual(st.query(0, 1), 9)
  assert.strictEqual(st.query(0, 3), 12)
  assert.strictEqual(st.query(0, 5), 15)
  assert.strictEqual(st.query(1, 1), 0)
  console.log(st.maxRight(0, preMax => preMax <= 12))
}
