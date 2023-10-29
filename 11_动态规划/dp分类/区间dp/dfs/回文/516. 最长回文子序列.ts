/** 最长回文子序列 (LPS). */
function longestPalindromeSubsequence(arr: ArrayLike<unknown>): number {
  const n = arr.length
  const dp = new Uint16Array(n * n)
  for (let i = n - 1; ~i; i--) {
    dp[i * n + i] = 1
    for (let j = i + 1; j < n; j++) {
      if (arr[i] === arr[j]) {
        dp[i * n + j] = dp[(i + 1) * n + j - 1] + 2
      } else {
        dp[i * n + j] = Math.max(dp[(i + 1) * n + j], dp[i * n + j - 1])
      }
    }
  }
  return dp[n - 1]
}

export { longestPalindromeSubsequence }

if (require.main === module) {
  console.log(longestPalindromeSubsequence('bbbab'))
  // https://leetcode.cn/problems/longest-palindromic-subsequence/
  // eslint-disable-next-line no-inner-declarations
  function longestPalindromeSubseq(s: string): number {
    return longestPalindromeSubsequence(s)
  }
}
