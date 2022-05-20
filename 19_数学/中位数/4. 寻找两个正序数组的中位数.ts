/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number}
 * 给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的 一半 。
 * 你能设计一个时间复杂度为 O(log (m+n)) 的算法解决此问题吗
 * @description 找出两个正序数组的中位数等价于找出两个正序数组中的第k小数
 */
function findMedianSortedArrays(nums1: number[], nums2: number[]): number {
  const n = nums1.length + nums2.length
  return (findK(nums1, nums2, n >> 1) + findK(nums1, nums2, (n - 1) >> 1)) / 2

  /**
   * @returns 寻找两个数组中第k小的数 k从0开始
   * @description log(m+n)
   */
  function findK(nums1: number[], nums2: number[], k: number): number {
    if (nums1.length === 0) return nums2[k]
    if (nums2.length === 0) return nums1[k]
    const i1 = nums1.length >> 1
    const i2 = nums2.length >> 1
    const m1 = nums1[i1]
    const m2 = nums2[i2]

    if (i1 + i2 < k) {
      // 如果 num1 的一半 大于nums2的一半 那么 nums2 的前半部分不包含第k小的数候选
      if (m1 > m2) return findK(nums1, nums2.slice(i2 + 1), k - (i2 + 1))
      else return findK(nums1.slice(i1 + 1), nums2, k - (i1 + 1))
    } else {
      // 如果 num1 的一半 大于nums2的一半 那么 nums1 的后半部分不包含第k小的数候选
      if (m1 > m2) return findK(nums1.slice(0, i1), nums2, k)
      else return findK(nums1, nums2.slice(0, i2), k)
    }
  }
}

console.log(findMedianSortedArrays([1, 2], [3, 4]))
export {}

/**
 * @param {number[]} arr1 - sorted integer array
 * @param {number[]} arr2 - sorted integer array
 * @returns {number}
 * O(m+n) 双指针 找出两个正序数组的中位数等价于找出两个正序数组中的第k小数
 * ps:有序数组 nums 的中位数为 (nums[n>>1]+nums[(n-1)>>1])/2
 */
function median(nums1: number[], nums2: number[]): number {
  const [n1, n2] = [nums1.length, nums2.length]
  const n = n1 + n2

  let [pre, cur] = [0, 0]
  let [i, j] = [0, 0]
  let step = (n >> 1) + 1

  while (step--) {
    pre = cur

    // 边界处理
    if (i === n1) {
      cur = nums2[j]
      j++
    } else if (j === n2) {
      cur = nums1[i]
      i++
    } else if (nums1[i] < nums2[j]) {
      cur = nums1[i]
      i++
    } else {
      cur = nums2[j]
      j++
    }
  }

  return (n & 1) === 0 ? (pre + cur) / 2 : cur
}
// console.log(median([1, 3], [2]))
// console.log(median([1, 3], [2, 4]))
