/**
 * 两个有序数组第k小的数.
 * !0 <= k < len(nums1) + len(nums2)
 */
function kthSmallestOfTwoSortedArrays(
  nums1: ArrayLike<number>,
  nums2: ArrayLike<number>,
  k: number
): number {
  let i1 = 0
  let i2 = 0
  let j1 = nums1.length
  let j2 = nums2.length
  while (true) {
    if (i1 === j1) return nums2[i2 + k]
    if (i2 === j2) return nums1[i1 + k]
    const mid1 = (j1 - i1) >>> 1
    const mid2 = (j2 - i2) >>> 1
    if (mid1 + mid2 < k) {
      if (nums1[i1 + mid1] < nums2[i2 + mid2]) {
        i1 += mid1 + 1
        k -= mid1 + 1
      } else {
        i2 += mid2 + 1
        k -= mid2 + 1
      }
    } else if (nums1[i1 + mid1] < nums2[i2 + mid2]) {
      j2 = i2 + mid2
    } else {
      j1 = i1 + mid1
    }
  }
}

export { kthSmallestOfTwoSortedArrays }

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  console.log(kthSmallestOfTwoSortedArrays([1, 3, 5, 7, 9], [2, 4, 6, 8, 10], 5)) // 5
}
