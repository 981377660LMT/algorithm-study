/**
 * @param {number[]} nums
 * @return {number[]}
 * 集合 s 包含从 1 到 n 的整数
 * 集合 丢失了一个数字 并且 有一个数字重复 。
 */
var findErrorNums = function (nums) {
  const res = []

  for (let i = 0; i < nums.length; i++) {
    const index = Math.abs(nums[i]) - 1
    if (nums[index] < 0) res.push(index + 1)
    else nums[index] *= -1
  }

  for (let i = 0; i < nums.length; i++) {
    if (nums[i] > 0) res.push(i + 1)
  }

  return res
}
