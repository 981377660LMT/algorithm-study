子序列默认不连续，子数组默认连续

最长重复子数组

```TS
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
```

最长公共子序列

```JS
 for (let i = 1; i < m + 1; i++) {
    for (let j = 1; j < n + 1; j++) {
      if (nums1[i - 1] === nums2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1
        res = Math.max(res, dp[i][j])
      }else{
        dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1])
      }
    }
  }
```
