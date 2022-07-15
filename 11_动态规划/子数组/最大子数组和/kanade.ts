/**
 * @description
 * 求最大/最小子数组和
 */
function kanade(nums: number[], getMax = true): number {
  if (nums.length === 0) return 0
  if (nums.length === 1) return nums[0]

  let res = getMax ? -Infinity : Infinity
  let dp = 0
  for (const num of nums) {
    if (getMax) {
      dp = Math.max(dp + num, num)
      res = Math.max(res, dp)
    } else {
      dp = Math.min(dp + num, num)
      res = Math.min(res, dp)
    }
  }

  return res
}

if (require.main === module) {
  console.log(kanade([-2, 1, -3, 4, -1, 2, 1, -5, 4]))
}

export { kanade }
