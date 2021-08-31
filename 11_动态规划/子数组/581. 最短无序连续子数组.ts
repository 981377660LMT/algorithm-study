/**
 * @param {number[]} nums
 * @return {number}
 * @description
 * 如果最右端的一部分已经排好序，这部分的每个数都比它左边的最大值要大
 * 如果最左端的一部分排好序，这每个数都比它右边的最小值小
 * 左往右遍历，如果i位置上的数比它左边部分最大值小，则这个数肯定要排序， 就这样找到右端不用排序的部分
 */
const findUnsortedSubarray = function (nums: number[]): number {
  const n = nums.length
  if (n <= 1) return 0

  let low = 0
  let high = 0
  let max = nums[0]
  let min = nums[n - 1]

  // high的右边可以不排序
  for (let i = 0; i < n; i++) {
    if (nums[i] >= max) {
      max = nums[i]
    } else {
      high = i
    }
  }

  // low的左边可以不排序
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
