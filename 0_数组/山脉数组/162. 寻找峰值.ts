// 你可以实现时间复杂度为 O(logN) 的解决方案吗？(暗示二分)
// 你可以假设 nums[-1] = nums[n] = -∞ 。
// 注意不能:
// arr.unshift(-Infinity)
// arr.push(-Infinity)

// 二分查找:考虑到开始导数大于0，最后导数小于0，因此如果mid导数大于0，则在右边，导数小于0则在左边
// 如果nums[i] > nums[i+1]，则在i之前一定存在峰值元素
// 如果nums[i] < nums[i+1]，则在i+1之后一定存在峰值元素
const findPeakElement = (nums: number[]) => {
  const n = nums.length
  let [l, r] = [0, n - 1]

  while (l <= r) {
    const mid = (l + r) >> 1
    const midVal = nums[mid]
    if ((midVal > nums[mid - 1] || mid === 0) && (midVal > nums[mid + 1] || mid === n - 1)) {
      return mid
    } else if (midVal > nums[mid - 1] || mid === 0) {
      l = mid + 1
    } else {
      r = mid - 1
    }
  }

  return l
}

console.log(findPeakElement([1, 2, 1, 3, 5, 6, 4]))
console.log(findPeakElement([1]))

export {}
