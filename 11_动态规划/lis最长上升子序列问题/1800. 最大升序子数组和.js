// 注意子数组与子序列区别
/**
 * @param {number[]} nums  正整数组成的数组 nums
 * @return {number}
 * 遍历即可
 */
function maxAscendingSum(nums) {
  if (nums.length === 1) return nums[0]
  let sum = nums[0]
  let res = nums[0]

  for (let i = 1; i < nums.length; i++) {
    if (nums[i] <= nums[i - 1]) sum = nums[i]
    else sum += nums[i]

    res = Math.max(res, sum)
  }

  return res
}

console.log(maxAscendingSum([10, 20, 30, 5, 10, 50]))
// 输出：65
// 解释：[5,10,50] 是元素和最大的升序子数组，最大元素和为 65 。
