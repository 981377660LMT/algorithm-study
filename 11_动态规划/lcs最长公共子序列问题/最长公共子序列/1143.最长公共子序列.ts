/**
 * LCS模板.Naive的实现.
 * {@link https://leetcode.cn/problems/longest-common-subsequence/}
 */
function longestCommonSubsequence1<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>): number {
  const n1 = arr1.length
  const n2 = arr2.length
  const dp: Uint32Array[] = Array(n1 + 1)
  for (let i = 0; i < n1 + 1; i++) dp[i] = new Uint32Array(n2 + 1)
  for (let i = 1; i < n1 + 1; i++) {
    const curRow = dp[i]
    const preRow = dp[i - 1]
    for (let j = 1; j < n2 + 1; j++) {
      if (arr1[i - 1] === arr2[j - 1]) {
        curRow[j] = preRow[j - 1] + 1
      } else {
        curRow[j] = Math.max(preRow[j], curRow[j - 1])
      }
    }
  }
  return dp[n1][n2]
}

/**
 * LCS滚动数组优化.
 * {@link https://leetcode.cn/problems/longest-common-subsequence/}
 */
function longestCommonSubsequence2<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>): number {
  const n1 = arr1.length
  const n2 = arr2.length
  const dp: Uint32Array[] = Array(n1 + 1)
  for (let i = 0; i < n1 + 1; i++) dp[i] = new Uint32Array(n2 + 1)
  for (let i = 1; i < n1 + 1; i++) {
    const curRow = dp[i]
    const preRow = dp[i - 1]
    for (let j = 1; j < n2 + 1; j++) {
      if (arr1[i - 1] === arr2[j - 1]) {
        curRow[j] = preRow[j - 1] + 1
      } else {
        curRow[j] = Math.max(preRow[j], curRow[j - 1])
      }
    }
  }
  return dp[n1][n2]
}

/**
 * 位运算加速LCS(最长公共子序列).
 * @complexity O(nm/w)
 */
function longestCommonSubsequence3<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>): number {
  const n1 = arr1.length
  const n2 = arr2.length
}

/**
 * O(n+m+(编辑距离)^2)求LCS，适用于相似度较高的情况.
 */
function longestCommonSubsequence4<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>): number {
  const n1 = arr1.length
  const n2 = arr2.length
}

export { longestCommonSubsequence3 as longestCommonSubsequence }

if (require.main === module) {
  console.log(longestCommonSubsequence1('abcde', 'ace'))
  console.log(longestCommonSubsequence1('abca', 'acba'))
}
