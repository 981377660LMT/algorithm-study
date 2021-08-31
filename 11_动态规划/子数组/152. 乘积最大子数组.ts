/**
 * @param {number[]} nums
 * @return {number}
 * @description 有点像 53. 最大子序和
 * 我们只要记录前i的最小值, 和最大值
 */
const maxProduct = function (nums: number[]): number {
  let res = nums[0]
  let preMax = nums[0]
  let preMin = nums[0]

  for (let i = 1; i < nums.length; i++) {
    const cur = nums[i]
    const curMax = Math.max(preMax * cur, preMin * cur, cur)
    const curMin = Math.min(preMax * cur, preMin * cur, cur)
    res = Math.max(res, curMax)
    preMax = curMax
    preMin = curMin
  }

  return res
}

console.log(maxProduct([2, 3, -2, 4]))
