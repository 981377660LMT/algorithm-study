/**
 * @param {number[]} nums
 * @return {number}
 题中递增排序数组为 0~n 中缺失一个数字。
 排序数组中的搜索问题，首先想到 二分法 解决
 */
const missingNumber = function (nums) {
  let l = 0
  let r = nums.length - 1
  while (l <= r) {
    const mid = (l + r) >> 1
    if (nums[mid] === mid) l = mid + 1
    else r = mid - 1
  }

  return l
}

console.log(missingNumber([0, 1, 2, 3, 4, 5, 6, 7, 9]))
// 8
