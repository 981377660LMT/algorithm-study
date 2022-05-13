import assert from 'assert'

type MergeFunc = (a: number, b: number) => number

/**
 * @summary st表适用于可重复贡献问题/RMQ静态区间最值查询
 * @description 可重复贡献问题是指对于运算`opt`，满足`x opt x = x`，则对应的区间询问就是一个可重复贡献问题。
 *
 * 例如，最大值有 `max(a,a) = a`，gcd有 `gcd(a,a) = a` ，所以RMQ和区间GCD就是一个可重复贡献问题。
 * 像区间和就不具有这个性质，如果求区间和的时候采用的预处理区间重叠了，则会导致重曼部分被计算两次，这是我们所不愿意看到的。
 * 另外，opt还必须满足结合律才能使用ST表求解。
 *
 * @see {@link https://oi-wiki.org/ds/sparse-table/}
 */
class SparseTable {
  private readonly size: number
  private readonly mergeFunc: MergeFunc
  private readonly dp: number[][]

  constructor(nums: number[], mergeFunc: MergeFunc) {
    const n = nums.length
    const upper = Math.ceil(Math.log2(n)) + 1

    this.size = n
    this.mergeFunc = mergeFunc
    this.dp = Array.from({ length: n }, () => Array(upper).fill(0))
    for (let i = 0; i < n; i++) this.dp[i][0] = nums[i]

    for (let j = 1; j < upper; j++) {
      for (let i = 0; i < n; i++) {
        if (i + (1 << (j - 1)) >= n) break
        this.dp[i][j] = this.mergeFunc(this.dp[i][j - 1], this.dp[i + (1 << (j - 1))][j - 1])
      }
    }
  }

  /**
   * @returns [`left`,`right`] 闭区间的贡献值
   */
  query(left: number, right: number): number {
    this.checkRange(left, right)
    const k = Math.floor(Math.log2(right - left + 1))
    return this.mergeFunc(this.dp[left][k], this.dp[right - (1 << k) + 1][k])
  }

  private checkRange(left: number, right: number): void {
    if (0 <= left && left <= right && right < this.size) return
    throw new RangeError(`invalid range [${left}, ${right}]`)
  }
}

if (require.main === module) {
  const st = new SparseTable([9, 12, 3, 7, 15], (a, b) => Math.max(a, b))
  assert.strictEqual(st.query(0, 0), 9)
  assert.strictEqual(st.query(0, 2), 12)
  assert.strictEqual(st.query(0, 4), 15)
}

export { SparseTable }
