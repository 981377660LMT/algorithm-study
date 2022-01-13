/**
 * @description
 * 求最大/最小子数组和
 */
function kanade(nums: number[], getMax = true): number {
  if (nums.length === 0) return 0
  if (nums.length === 1) return nums[0]

  let res = getMax ? -Infinity : Infinity
  let sum = 0
  for (const num of nums) {
    if (getMax) {
      sum = Math.max(sum + num, num)
      res = Math.max(res, sum)
    } else {
      sum = Math.min(sum + num, num)
      res = Math.min(res, sum)
    }
  }

  return res
}

// console.log(maxSubArray([-2, 1, -3, 4, -1, 2, 1, -5, 4]))

export { kanade }
