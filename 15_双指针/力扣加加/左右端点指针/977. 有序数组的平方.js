/**
 * @param {number[]} nums
 * @return {number[]}
 * 给你一个按 `非递减顺序` 排序的整数数组 nums
 * 回 每个数字的平方 组成的新数组，要求也按 非递减顺序 排序。
 * 负数平方之后可能成为最大数:从大排到小
 */
function sortedSquares(nums) {
  const res = Array(nums.length).fill(0)
  let l = 0
  let r = nums.length - 1

  for (let i = nums.length - 1; i >= 0; i--) {
    if (nums[l] ** 2 > nums[r] ** 2) {
      res[i] = nums[l] ** 2
      l++
    } else {
      res[i] = nums[r] ** 2
      r--
    }
  }

  return res
}

console.log(sortedSquares([-4, -1, 0, 3, 10]))
