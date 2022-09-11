/* eslint-disable func-names */
/* eslint-disable no-param-reassign */

// 给定两个排序后的数组 A 和 B，其中 A 的末端有足够的缓冲空间容纳 B。
// 编写一个方法，将 B 合并入 A 并排序。
/**
 * @return {void} Do not return anything, modify nums1 in-place instead.
 * @description 从后向前排,nums2用完即结束
 */
function merge(nums1, m, nums2, n) {
  let i = m - 1
  let j = n - 1
  let k = m + n - 1

  // nums2用完即结束
  while (j >= 0) {
    // 注意nums1为空的情况
    if (nums1[i] > nums2[j]) {
      nums1[k] = nums1[i]
      i--
    } else {
      nums1[k] = nums2[j]
      j--
    }

    k--
  }

  return nums1
}

console.log(merge([1, 2, 3, 0, 0, 0], 3, [2, 5, 6], 3))
console.log(merge([], 0, [1], 1))
// 输出：[1,2,2,3,5,6]
