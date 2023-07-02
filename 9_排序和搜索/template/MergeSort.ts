/* eslint-disable no-inner-declarations */

function mergeSort<T>(arr: ArrayLike<T>, compareFn: (a: T, b: T) => number): T[] {
  if (!arr.length) return []
  if (arr.length === 1) return [arr[0]]

  const copy: T[][] = Array(arr.length)
  for (let i = 0; i < arr.length; i++) copy[i] = [arr[i]]

  // 分治的迭代实现
  let n = arr.length
  while (n > 1) {
    const mid = (n + 1) >>> 1
    for (let i = 0; i < mid; i++) {
      if (((i << 1) | 1) ^ n) {
        copy[i] = merge(copy[i << 1], copy[(i << 1) | 1])
      } else {
        copy[i] = copy[i << 1]
      }
    }
    n = mid
  }

  return copy[0]

  function merge(nums1: ArrayLike<T>, nums2: ArrayLike<T>): T[] {
    const res = Array<T>(nums1.length + nums2.length)
    let i = 0
    let j = 0
    let ptr = 0
    while (i < nums1.length && j < nums2.length) {
      if (compareFn(nums1[i], nums2[j]) < 0) {
        res[ptr++] = nums1[i++]
      } else {
        res[ptr++] = nums2[j++]
      }
    }

    while (i < nums1.length) res[ptr++] = nums1[i++]
    while (j < nums2.length) res[ptr++] = nums2[j++]
    return res
  }
}

export {}

if (require.main === module) {
  const arr = [4, 2, 100, 99, 10000, -1, 99, 2]
  console.log(mergeSort(arr, (a, b) => a - b))

  // https://leetcode.cn/problems/sort-an-array/
  function sortArray(nums: number[]): number[] {
    return mergeSort(nums, (a, b) => a - b)
  }
}
