const INF = 2e15

/**
 * @param {number[]} nums
 * @param {number} x
 * @return {number}
 * @description 每一次操作时，你应当移除数组 nums 最左边或最右边的元素，然后从 x 中减去该元素的值。请注意，需要 修改 数组以供接下来的操作使用。
 * 返回 最小操作数
 * @summary 求和为定值sum(nums)-x的最长子数组: 滑动窗口
 */
function minOperations(nums: number[], x: number): number {
  const target = nums.reduce((pre, cur) => pre + cur, 0) - x
  let res = -INF
  let curSum = 0
  let left = 0

  for (let right = 0; right < nums.length; right++) {
    curSum += nums[right]

    while (left <= right && curSum > target) {
      curSum -= nums[left]
      left++
    }

    if (curSum === target) {
      res = Math.max(res, right - left + 1)
    }
  }

  return res === -INF ? -1 : nums.length - res
}

console.log(
  minOperations(
    [
      8828, 9581, 49, 9818, 9974, 9869, 9991, 10000, 10000, 10000, 9999, 9993, 9904, 8819, 1231,
      6309
    ],
    134365
  )
)
