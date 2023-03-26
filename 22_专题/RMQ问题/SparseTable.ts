import assert from 'assert'

type MergeFunc<S> = (a: S, b: S) => S

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
  private readonly _size: number
  private readonly _mergeFunc: MergeFunc<S>
  private readonly _dp: S[][]

  constructor(nums: ArrayLike<S>, mergeFunc: MergeFunc<S>) {
    const n = nums.length
    const size = 32 - Math.clz32(n)

    this._size = n
    this._mergeFunc = mergeFunc
    // !dp[i][j]表示区间[j,j+2**i-1]的最大值
    this._dp = Array.from({ length: size }, () => Array(n).fill(0))
    for (let i = 0; i < n; i++) this._dp[0][i] = nums[i]

    for (let i = 1; i < size; i++) {
      for (let j = 0; j < n - (1 << i) + 1; j++) {
        this._dp[i][j] = mergeFunc(this._dp[i - 1][j], this._dp[i - 1][j + (1 << (i - 1))])
      }
    }
  }

  /**
   * @returns [`left`,`right`] 闭区间的贡献值
   * @param left 0 <= left <= right < nums.length
   * @param right 0 <= left <= right < nums.length
   */
  query(left: number, right: number): S {
    // this._checkBoundsBeginEnd(left, right)
    const k = 32 - Math.clz32(right - left + 1) - 1
    return this._mergeFunc(this._dp[k][left], this._dp[k][right - (1 << k) + 1])
  }

  private _checkBoundsBeginEnd(begin: number, end: number): void {
    if (begin >= 0 && begin <= end && end < this._size) return
    throw new RangeError(`invalid range [${begin}, ${end}]`)
  }
}

if (require.main === module) {
  const st = new SparseTable([9, 12, 3, 7, 15], (a, b) => Math.max(a, b))
  assert.strictEqual(st.query(0, 0), 9)
  assert.strictEqual(st.query(0, 2), 12)
  assert.strictEqual(st.query(0, 4), 15)
}

export { SparseTable }
