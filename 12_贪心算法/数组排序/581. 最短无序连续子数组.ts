/**
 * @param {number[]} nums
 * @return {number}
 * @description
给你一个整数数组 nums ，你需要找出一个 连续子数组 ，如果对这个子数组进行升序排序，那么整个数组都会变为升序排序。就这样找到右端不用排序的部分
 */
const findUnsortedSubarray = function (nums: number[]): number {
  const n = nums.length
  if (n <= 1) return 0

  let low = 0
  let high = 0
  let max = nums[0]
  let min = nums[n - 1]

  // 从左到右找high high的右边可以不排序
  for (let i = 0; i < n; i++) {
    if (nums[i] >= max) {
      max = nums[i]
    } else {
      high = i
    }
  }

  // 从右到左找low low的左边可以不排序
  for (let i = n - 1; i >= 0; i--) {
    if (nums[i] <= min) {
      min = nums[i]
    } else {
      low = i
    }
  }

  return high === low ? 0 : high - low + 1
}

console.log(findUnsortedSubarray([2, 6, 4, 8, 10, 9, 15]))
console.log(findUnsortedSubarray([2, 1]))

export {}
