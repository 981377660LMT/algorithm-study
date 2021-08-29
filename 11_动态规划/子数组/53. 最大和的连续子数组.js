/**
 * @param {number[]} nums
 * @return {number}
 */
var maxSubArray = function (nums) {
  if (nums.length === 1) return nums[0]
  let sum = 0
  let res = -Infinity
  for (const num of nums) {
    sum = Math.max(sum + num, num)
    res = Math.max(res, sum)
  }
  return res
}
