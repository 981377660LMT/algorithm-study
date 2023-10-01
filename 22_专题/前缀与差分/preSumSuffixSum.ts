/* eslint-disable no-inner-declarations */

// MonoidPresum/PresumMonoid/PreSumSuffixSum

class PreSumSuffixSum<E> {
  private readonly _preSum: E[]
  private readonly _suffixSum: E[]
  private readonly _e: () => E

  constructor(arr: ArrayLike<E>, e: () => E, op: (a: E, b: E) => E) {
    const n = arr.length
    const preSum: E[] = Array(n + 1).fill(e())
    const suffixSum: E[] = Array(n + 1).fill(e())
    preSum[0] = e()
    suffixSum[n] = e()
    for (let i = 0; i < n; i++) {
      preSum[i + 1] = op(preSum[i], arr[i])
      suffixSum[n - i - 1] = op(suffixSum[n - i], arr[n - i - 1])
    }
    this._e = e
    this._preSum = preSum
    this._suffixSum = suffixSum
  }

  /**
   * 查询前缀 `[0,end)` 的和.
   */
  preSum(end: number): E {
    if (end < 0) return this._e()
    if (end >= this._preSum.length) return this._preSum[this._preSum.length - 1]
    return this._preSum[end]
  }

  /**
   * 查询后缀 `[start,n)` 的和.
   */
  suffixSum(start: number): E {
    if (start < 0) return this._suffixSum[0]
    if (start >= this._suffixSum.length) return this._e()
    return this._suffixSum[start]
  }
}

export { PreSumSuffixSum }

if (require.main === module) {
  // https://leetcode.cn/problems/maximum-value-of-an-ordered-triplet-ii/description/
  function maximumTripletValue(nums: number[]): number {
    const P = new PreSumSuffixSum(nums, () => 0, Math.max)
    let res = 0
    for (let j = 1; j < nums.length - 1; j++) {
      const preMax = P.preSum(j)
      const sufMax = P.suffixSum(j + 1)
      res = Math.max(res, (preMax - nums[j]) * sufMax)
    }
    return res
  }
}
