1. 子序列默认不连续，子数组默认连续

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

2. 最长公共子序列是否存在低于 O(n^2) 的算法？
   https://www.zhihu.com/question/21974937
   diff 本质上就是一个 LCS，只不过 diff 可能会允许近似解
   git diff 了一个近万行的程序，但是实际上你可能只修改了不到一百行。
   **衡量一个 diff 算法的优劣指标不应该仅仅由长度 n 决定，还应该由差异数量 d 决定**

   `元素不同的 lcs 可以转化为 lis`
