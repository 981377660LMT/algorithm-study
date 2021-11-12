/**
 * @param {number[]} nums   可能存在 重复 元素值的数组 nums
 * @return {number}
 * 最小元素特点:左侧元素都大于数组第一个元素,右侧元素都小于数组第一个元素
 * 平均时间复杂度为 O(logn)，而在最坏情况下，
 * 如果数组中的元素完全相同，那么 while 循环就需要执行 n次，
 * 每次忽略区间的右端点，时间复杂度为 O(n)。
 */
const findMin = function (nums: number[]): number {
  let l = 0
  let r = nums.length - 1

  // 循环外返回，不能取等号
  while (l < r) {
    const mid = (l + r) >> 1
    if (nums[mid] === nums[r]) {
      // 重复多了这一句
      // 当中间值等于右边，将右边界移过来，因为左边可能还有相等的值
      r--
    } else if (nums[mid] > nums[r]) {
      l = mid + 1
    } else if (nums[mid] < nums[r]) {
      r = mid
    }
  }

  return nums[l]
}

console.log(findMin([4, 5, 6, 7, 0, 1, 2]))
console.log(findMin([1, 1]))
console.log(findMin([3, 1, 3]))

export {}
