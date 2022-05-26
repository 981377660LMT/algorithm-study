// 给你一个长度为 n 的整数数组，请你判断在 最多 改变 1 个元素的情况下，该数组能否变成一个非递减数列。
/**
 * @param {number[]} nums
 * @return {boolean}
 * @link https://leetcode-cn.com/problems/non-decreasing-array/solution/3-zhang-dong-tu-bang-zhu-ni-li-jie-zhe-d-06gi/
 * @summary
 * 当 nums[i] 破坏了数组的单调递增时，即 nums[i] < nums[i - 1]  时，为了让数组有序，我们发现一个规律
 * 1. 当 i = 1 ，那么修改 num[i- 1] ，不要动 nums[i]
 * 2.当 i > 1 时，我们应该优先考虑把 nums[i - 1] 调小到 >= nums[i - 2] 并且 <= nums[i]。
 * 3. 当 i > 1 且 nums[i] < nums[i - 2] 时，我们无法调整 nums[i - 1] ，我们只能调整 nums[i] 到 nums[i - 1] 。
 * 即：尽量调前面，实在不行调后面
 */
var checkPossibility = function (nums) {
  let res = 0
  const n = nums.length
  for (let i = 1; i < n; i++) {
    if (nums[i] < nums[i - 1]) {
      res++
      if (i === 1 || nums[i] >= nums[i - 2]) {
        nums[i - 1] = nums[i]
      } else {
        nums[i] = nums[i - 1]
      }
    }
  }

  return res <= 1
}
