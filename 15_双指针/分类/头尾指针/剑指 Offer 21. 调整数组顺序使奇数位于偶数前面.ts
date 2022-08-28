/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable no-param-reassign */
/* eslint-disable semi-style */

/**
 * @param {number[]} nums
 * @return {number[]}
 * 使得所有奇数位于数组的前半部分，所有偶数位于数组的后半部分
 * !头尾双指针 头找偶 尾找奇
 */
function exchange(nums: number[]): number[] {
  if (nums.length === 1) {
    return nums
  }

  let left = 0
  let right = nums.length - 1
  while (left < right) {
    while (left < right && nums[left] & 1) {
      left++
    }

    while (left < right && !(nums[right] & 1)) {
      right--
    }

    ;[nums[left], nums[right]] = [nums[right], nums[left]]
  }

  return nums
}
