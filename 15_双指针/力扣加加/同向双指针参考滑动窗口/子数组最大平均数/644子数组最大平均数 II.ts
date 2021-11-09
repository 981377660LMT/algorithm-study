/**
 * @param {number[]} nums  -104 <= nums[i] <= 10
 * @param {number} k  1 <= k <= n <= 104
 * @return {number}
 * 找出 长度大于等于 k 且含最大平均值的连续子数组 并输出这个最大平均值
 * 任何计算误差小于 10-5 的结果都将被视为正确答案
 * 答案误差不超过1e-5是使用二分搜索的提示
 * 二分 nlog(max-min)
 * @description 二分答案法
 * @summary 二分 子数组 前缀和
 */
const findMaxAverage = function (nums: number[], k: number): number {
  let l = Math.min(...nums)
  let r = Math.max(...nums)

  while (r - l >= 1e-5) {
    const mid = (l + r) / 2
    if (check(mid)) l = mid
    else r = mid
  }

  return l

  // 存在长度不小于k的子数组平均数大于等于average 使用前缀和求子数组和
  function check(average: number): boolean {
    const preSum = Array<number>(nums.length + 1).fill(0)
    for (let i = 1; i < preSum.length; i++) {
      preSum[i] = preSum[i - 1] + nums[i - 1] - average // 平移
    }

    let preMin = 0
    for (let i = k; i < preSum.length; i++) {
      if (preSum[i] - preMin >= 0) return true
      preMin = Math.min(preMin, preSum[i - k + 1])
    }

    return false
  }
}

console.log(findMaxAverage([1, 12, -5, -6, 50, 3], 4))
