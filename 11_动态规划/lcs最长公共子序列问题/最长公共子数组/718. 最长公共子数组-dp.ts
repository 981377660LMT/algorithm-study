/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number}
 * @description  dp[i][j] 为 以 A[i], B[j] 结尾的两个数组中公共的、长度最长的子数组的长度
 */
const findLength = function (nums1: number[], nums2: number[]): number {
  let res = 0
  const m = nums1.length
  const n = nums2.length
  const dp = Array.from<number, number[]>({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i < m + 1; i++) {
    for (let j = 1; j < n + 1; j++) {
      if (nums1[i - 1] === nums2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1
        res = Math.max(res, dp[i][j])
      }
    }
  }

  return res
}

console.log(findLength([1, 2, 3, 2, 1], [3, 2, 1, 4, 7]))

export default 1
