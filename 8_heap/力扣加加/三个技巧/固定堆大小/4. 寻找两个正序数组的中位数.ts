/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number}
 * 给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的 一半 。
 * 你能设计一个时间复杂度为 O(log (m+n)) 的算法解决此问题吗
 * @description 找出两个正序数组的中位数等价于找出两个正序数组中的第k小数
 */
var findMedianSortedArrays = function (nums1: number[], nums2: number[]): number {
  // k 从 0 开始
  // 每次寻找二分索引
  const findK = (nums1: number[], nums2: number[], k: number): number => {
    if (!nums1.length) return nums2[k]
    if (!nums2.length) return nums1[k]
    const i1 = nums1.length >> 1
    const i2 = nums2.length >> 1
    const m1 = nums1[i1]
    const m2 = nums2[i2]
    console.log(i1, i2, nums1, nums2, k)
    if (i1 + i2 < k) {
      // 如果 num1 的一半 大于nums2的一半 那么 nums2 的前半部分不包含第k小的数候选
      if (m1 > m2) return findK(nums1, nums2.slice(i2 + 1), k - i2 - 1)
      else return findK(nums1.slice(i1 + 1), nums2, k - i1 - 1)
    } else {
      // 如果 num1 的一半 大于nums2的一半 那么 nums1 的后半部分不包含第k小的数候选
      if (m1 > m2) return findK(nums1.slice(0, i1), nums2, k)
      else return findK(nums1, nums2.slice(0, i2), k)
    }
  }
  // return findK([1, 2, 3], [4, 5, 6], 3)

  const len = nums1.length + nums2.length
  if (len % 2 === 1) return findK(nums1, nums2, len >> 1)
  else return (findK(nums1, nums2, len / 2) + findK(nums1, nums2, len / 2 - 1)) / 2
}

console.log(findMedianSortedArrays([1, 2], [3, 4]))
export {}
