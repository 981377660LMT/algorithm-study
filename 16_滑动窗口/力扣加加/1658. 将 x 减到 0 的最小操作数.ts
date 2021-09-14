/**
 * @param {number[]} nums
 * @param {number} x
 * @return {number}
 * @description 每一次操作时，你应当移除数组 nums 最左边或最右边的元素，然后从 x 中减去该元素的值。请注意，需要 修改 数组以供接下来的操作使用。
 * 返回 最小操作数
 * @summary 求和为定值sum(nums)-x的最长子数组: 滑动窗口
 */
const minOperations = function (nums: number[], x: number): number {
  let l = 0
  let res = 0
  let sum = 0
  const target = nums.reduce((pre, cur) => pre + cur, 0) - x

  // 注意这里
  if (target === 0) return nums.length

  for (let r = 0; r < nums.length; r++) {
    sum += nums[r]
    while (sum > target) {
      l++
      sum -= nums[l - 1]
    }
    if (sum === target) {
      res = Math.max(res, r - l + 1)
    }
  }

  return res === 0 ? -1 : nums.length - res
}

console.log(
  minOperations(
    [
      8828, 9581, 49, 9818, 9974, 9869, 9991, 10000, 10000, 10000, 9999, 9993, 9904, 8819, 1231,
      6309,
    ],
    134365
  )
)
