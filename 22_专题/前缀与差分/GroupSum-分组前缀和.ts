/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */

/**
 * 模分组前缀和/同余前缀和.
 * @param arr 给定的数组.
 * @param mod 模数.
 * @returns 求下标在[start,end)范围内, 模mod为key的元素的和.
 */
function groupPresum(arr: ArrayLike<number>, mod: number): (start: number, end: number, key: number) => number {
  const preSum = Array(arr.length + mod).fill(0)
  for (let i = 0; i < arr.length; i++) {
    preSum[i + mod] = preSum[i] + arr[i]
  }

  const cal = (r: number, k: number): number => {
    if (r % mod <= k) {
      return preSum[Math.floor(r / mod) * mod + k]
    }
    return preSum[Math.floor((r + mod - 1) / mod) * mod + k]
  }

  const query = (start: number, end: number, key: number): number => {
    if (start >= end) return 0
    key %= mod
    return cal(end, key) - cal(start, key)
  }

  return query
}

export { groupPresum }

if (require.main === module) {
  const S = groupPresum([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 3)
  console.log(S(0, 10, 0))
  console.log(S(0, 10, 1))
  console.log(S(0, 10, 2))

  // LC1664 https://leetcode.cn/problems/ways-to-make-a-fair-array/
  // 1664. 生成平衡数组的方案数
  function waysToMakeFair(nums: number[]): number {
    let res = 0
    const sum = groupPresum(nums, 2)
    for (let i = 0; i < nums.length; i++) {
      const leftOddSum = sum(0, i, 1)
      const leftEvenSum = sum(0, i, 0)
      const rightOddSum = sum(i + 1, nums.length, 1)
      const rightEvenSum = sum(i + 1, nums.length, 0)
      res += +(leftOddSum + rightEvenSum === leftEvenSum + rightOddSum)
    }
    return res
  }
}
