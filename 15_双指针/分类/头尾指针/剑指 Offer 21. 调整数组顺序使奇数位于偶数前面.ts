/**
 * @param {number[]} nums
 * @return {number[]}
 * 使得所有奇数位于数组的前半部分，所有偶数位于数组的后半部分
 * 奇偶双指针 头找偶 尾找奇
 */
const exchange = function (nums: number[]): number[] {
  if (nums.length === 1) return nums
  let i = 0
  let j = nums.length - 1
  while (i < j) {
    while (i < j && nums[i] & 1) i++
    while (i < j && !(nums[j] & 1)) j--
    ;[nums[i], nums[j]] = [nums[j], nums[i]]
  }
  return nums
}
