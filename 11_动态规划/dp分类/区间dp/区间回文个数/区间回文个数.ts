/**
 * 统计区间回文串个数.
 */
function rangePalindrome(arr: ArrayLike<unknown>): (start: number, end: number) => number {
  const n = arr.length
  const dp = new Int32Array(n * n)
  for (let i = 0; i < n; i++) {
    dp[i * n + i] = 1
    if (i + 1 < n && arr[i] === arr[i + 1]) {
      dp[i * n + i + 1] = 1
    }
  }

  for (let i = n - 3; i >= 0; i--) {
    for (let j = i + 2; j < n; j++) {
      if (arr[i] === arr[j]) {
        dp[i * n + j] = dp[(i + 1) * n + j - 1]
      }
    }
  }

  for (let i = n - 2; i >= 0; i--) {
    for (let j = i + 1; j < n; j++) {
      dp[i * n + j] += dp[i * n + j - 1] + dp[(i + 1) * n + j] - dp[(i + 1) * n + j - 1]
    }
  }

  return (start: number, end: number): number => {
    if (start < 0) start = 0
    if (end > n) end = n
    if (start >= end) return 0
    return dp[start * n + end - 1]
  }
}

export { rangePalindrome }

if (require.main === module) {
  const R = rangePalindrome('abab')
  console.log(R(0, 3))
}
