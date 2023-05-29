/* eslint-disable no-param-reassign */
import assert from 'assert'

/**
 * st表适用于可重复贡献问题/RMQ静态区间最值查询。
 * 可重复贡献问题是指对于运算`opt`，满足`x opt x = x`，则对应的区间询问就是一个可重复贡献问题。
 * !幂等性
 *
 * 例如，最大值有 `max(a,a) = a`，gcd有 `gcd(a,a) = a` ，所以RMQ和区间GCD就是一个可重复贡献问题。
 * 像区间和就不具有这个性质，如果求区间和的时候采用的预处理区间重叠了，
 * 则会导致重曼部分被计算两次，这是我们所不愿意看到的。
 * 另外，opt还必须满足结合律才能使用ST表求解。
 *
 * @see {@link https://oi-wiki.org/ds/sparse-table/}
 * @important cpu-cache optimized dp[bit][n]
 */
class SparseTable<S = number> {
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S
  private readonly _dp: S[][]
  private readonly _n: number

  constructor(nums: ArrayLike<S>, e: () => S, op: (a: S, b: S) => S) {
    const n = nums.length
    const size = 32 - Math.clz32(n)
    this._n = n
    this._e = e
    this._op = op
    this._dp = Array(size) // !dp[i][j]表示区间[j,j+2**i)的贡献值
    for (let i = 0; i < size; i++) this._dp[i] = Array(n)
    for (let i = 0; i < n; i++) this._dp[0][i] = nums[i]
    for (let i = 1; i < size; i++) {
      for (let j = 0; j < n - (1 << i) + 1; j++) {
        this._dp[i][j] = op(this._dp[i - 1][j], this._dp[i - 1][j + (1 << (i - 1))])
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
    return this._op(this._dp[k][start], this._dp[k][end - (1 << k)])
  }

  toString(): string {
    return `SparseTable{${this._dp[0]}}`
  }
}

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
  console.log(st.toString())
}

export { SparseTable }
