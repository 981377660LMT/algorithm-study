// 你可以实现时间复杂度为 O(logN) 的解决方案吗？(暗示二分)
// 你可以假设 nums[-1] = nums[n] = -∞ 。
// 线性查找:等价于找到第一个满足nums[i]>nums[i+1]的i
// 二分查找:考虑到开始导数大于0，最后导数小于0，因此如果mid导数大于0，则在右边，导数小于0则在左边
// 如果nums[i] > nums[i+1]，则在i之前一定存在峰值元素
// 如果nums[i] < nums[i+1]，则在i+1之后一定存在峰值元素
const findPeakElement = (nums: number[]) => {
  let [l, r] = [0, nums.length - 1]

  while (l < r) {
    const mid = (l + r) >> 1

    if (nums[mid] > nums[mid + 1]) {
      r = mid
    } else {
      l = mid + 1
    }
  }

  return l
}

console.log(findPeakElement([1, 2, 1, 3, 5, 6, 4]))

export {}
