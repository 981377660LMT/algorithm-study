/**
 * @param {number[]} nums
 * @return {number}
 * @description 对每一个数 带还是不带
 * “最大”“连续”，可以判断是一道动态规划
 * f(n)表示以第n个数为结束点的子数列的最大和
 * @summary
 * kanade 算法
 */
const maxSubArray = function (nums: number[]): number {
  if (nums.length === 1) return nums[0]
  let sum = 0
  let res = -Infinity
  for (const num of nums) {
    sum = Math.max(sum + num, num)
    res = Math.max(res, sum)
  }
  return res
}

// console.log(maxSubArray([-2, 1, -3, 4, -1, 2, 1, -5, 4]))

export { maxSubArray }
