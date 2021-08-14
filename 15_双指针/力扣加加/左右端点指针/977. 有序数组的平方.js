/**
 * @param {number[]} nums
 * @return {number[]}
 * 给你一个按 非递减顺序 排序的整数数组 nums
 * 回 每个数字的平方 组成的新数组，要求也按 非递减顺序 排序。
 */
var sortedSquares = function (nums) {
  const res = Array(nums.length)
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
