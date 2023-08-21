const INF = 2e15

/**
 * O(log(m+n))求两个正序数组的中位数.
 * @see {@link https://leetcode.cn/problems/median-of-two-sorted-arrays/solutions/258842/xun-zhao-liang-ge-you-xu-shu-zu-de-zhong-wei-s-114/}
 */
function findMedian(nums1: ArrayLike<number>, nums2: ArrayLike<number>): number {
  if (nums1.length > nums2.length) {
    const tmp = nums1
    nums1 = nums2
    nums2 = tmp
  }

  const len1 = nums1.length
  const len2 = nums2.length
  let left = 0
  let right = len1
  let max1 = 0 // 前一部分最大值
  let min2 = 0 // 后一部分最小值

  while (left <= right) {
    const i = (left + right) >>> 1 // 前一部分包含 nums1[0 .. i-1] 和 nums2[0 .. j-1]
    const j = ((len1 + len2 + 1) >>> 1) - i // 后一部分包含 nums1[i .. m-1] 和 nums2[j .. n-1]
    const a1 = i ? nums1[i - 1] : -INF
    const b1 = i < len1 ? nums1[i] : INF
    const a2 = j ? nums2[j - 1] : -INF
    const b2 = j < len2 ? nums2[j] : INF

    if (a1 <= b2) {
      max1 = a1 > a2 ? a1 : a2
      min2 = b1 < b2 ? b1 : b2
      left = i + 1
    } else {
      right = i - 1
    }
  }

  return (len1 + len2) & 1 ? max1 : (max1 + min2) / 2
}

export { findMedian }

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  function findMedianSortedArrays(nums1: number[], nums2: number[]): number {
    return findMedian(nums1, nums2)
  }
}
