/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number}
 * @summary
 * 无需考虑全部整体，而只需要考虑相邻两个数字即可
 * 相邻两个数字的大小关系有哪些?
 * 
 * q1 表示的是两个数组本身就已经递增了，你可以选择不交换。
   q2 表示的是两个数组必须进行一次交换，你可以选择交换 i 或者交换 i - 1。
 */
const minSwap = function (nums1: number[], nums2: number[]): number {
  const len = nums1.length
  // dp_swap[i]表示A[0:i+1],B[0:i+1]在交换A[i]与B[i]并保持数组递增的最小交换次数
  // dp_keep[i]表示A[0:i+1],B[0:i+1]在不交换A[i]与B[i]并保持数组递增的最小交换次数
  const swap = Array<number>(len).fill(Infinity)
  const noSwap = Array<number>(len).fill(Infinity)
  swap[0] = 1
  noSwap[0] = 0

  for (let i = 0; i < len; i++) {
    // 如果交换之前有序，则可以不交换
    if (nums1[i] > nums1[i - 1] && nums2[i] > nums2[i - 1]) {
      swap[i] = swap[i - 1] + 1
      noSwap[i] = noSwap[i - 1]
    }

    // 否则至少需要交换一次（交换当前项或者前一项）
    if (nums1[i] > nums2[i - 1] && nums2[i] > nums1[i - 1]) {
      swap[i] = Math.min(swap[i], noSwap[i - 1] + 1)
      noSwap[i] = Math.min(noSwap[i], swap[i - 1])
    }
  }

  return Math.min(swap[len - 1], noSwap[len - 1])
}

console.log(minSwap([1, 3, 5, 4], [1, 2, 3, 7]))

export {}
