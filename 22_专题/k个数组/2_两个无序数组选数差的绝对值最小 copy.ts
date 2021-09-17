/**
 *
 * @param arr1
 * @param arr2
 * @description 给你两个无序的非空数组 nums1 和 nums2，让你从每个数组中分别挑一个，使得二者差的绝对值最小。
 * 时间复杂度：$O(NlogN)$
 */
const minDiff = (arr1: number[], arr2: number[]): number => {
  arr1.sort((a, b) => a - b)
  arr2.sort((a, b) => a - b)
  let res = Infinity
  let l1 = 0
  let l2 = 0
  while (l1 < arr1.length && l2 < arr2.length) {
    res = Math.min(res, Math.abs(arr1[l1] - arr2[l2]))
    // 找最小值,小的向后移
    arr1[l1] < arr2[l2] ? l1++ : l2++
  }
  return res
}

console.log(minDiff([10, 30, 20, 30, 40], [25, 28, 7]))
export default 1
