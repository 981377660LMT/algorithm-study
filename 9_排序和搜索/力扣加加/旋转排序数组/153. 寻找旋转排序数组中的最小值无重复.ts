/**
 * @param {number[]} nums  元素值 互不相同 的数组 nums
 * @return {number}
 */
const findMin = function (nums: number[]): number {
  let l = 0
  let r = nums.length - 1

  while (l < r) {
    const mid = (l + r) >> 1
    if (nums[mid] > nums[r]) {
      l = mid + 1
    } else if (nums[mid] < nums[r]) {
      // 不能 mid-1 如果mid最小值就跳过了
      r = mid
    }
  }

  // 这里
  return nums[l]
}

console.log(findMin([3, 4, 5, 1, 2]))
