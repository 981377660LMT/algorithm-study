/**
 * @param {number[]} nums
 * @return {number}
 题中递增排序数组为 0~n 中缺失一个数字。
 排序数组中的搜索问题，首先想到 二分法 解决
 */
const missingNumber = function (nums) {
  let left = 0
  let right = nums.length - 1
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (nums[mid] === mid) left = mid + 1
    else right = mid - 1
  }

  return left
}

console.log(missingNumber([0, 1, 2, 3, 4, 5, 6, 7, 9]))
// 8
