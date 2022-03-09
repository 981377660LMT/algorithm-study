/**
 * @param {number[]} nums
 * @return {boolean}
 * 判断这个数组中是否存在长度为 3 的递增子序列。
 * 要求 O(n)的时间复杂度和 O(1)的空间复杂度
 * @summary
 * 三个变量，分别记录最小值，第二小值，第三小值。
 * 只要我们能够填满这三个变量就返回 true，否则返回 false。
 */
function increasingTriplet(nums: number[]): boolean {
  if (nums.length <= 2) return false

  let left = nums[0]
  let mid = Infinity
  for (const num of nums) {
    if (num <= left) {
      left = num
    } else if (num <= mid) {
      mid = num
    } else {
      return true
    }
  }

  return false
}

console.log(increasingTriplet([1, 2, 3, 4, 5]))
