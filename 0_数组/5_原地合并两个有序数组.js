/**
 * @param {number[]} nums1
 * @param {number} m
 * @param {number[]} nums2
 * @param {number} n
 * @return {void} Do not return anything, modify nums1 in-place instead.
 * @description 从后向前排,nums2用完即结束
 */
var merge = function (nums1, m, nums2, n) {
  let insertPostion = m + n - 1
  m--
  n--
  // nums2用完即结束
  while (n >= 0) {
    // 注意nums1为空的情况
    if (nums1[m] > nums2[n]) {
      nums1[insertPostion] = nums1[m]
      m--
    } else {
      nums1[insertPostion] = nums2[n]
      n--
    }

    insertPostion--
  }
  return nums1
}

console.log(merge([1, 2, 3, 0, 0, 0], 3, [2, 5, 6], 3))
console.log(merge([], 0, [1], 1))
// 输出：[1,2,2,3,5,6]
