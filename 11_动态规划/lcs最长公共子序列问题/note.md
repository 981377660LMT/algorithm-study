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

---

更快的 LCS

1. 相似度较高(超过 90％)时的线性算法:O(n+m+(编辑距离)^2)
   https://rsk0315.hatenablog.com/entry/2022/12/30/180216
   https://atcoder.jp/contests/dp/editorial/3793
   基因匹配可能有那种相似度极高的场景
   互联网领域感觉很少
2. 通用的位运算加速 O(nm/w)
   https://www.cnblogs.com/-Wallace-/p/bit-lcs.html
   https://wenku.baidu.com/view/ed99e4f77c1cfad6195fa776.html?_wkts_=1688222807315
   https://leetcode.cn/problems/longest-common-subsequence/solution/onmw-0ms-100tie-ge-mo-ban-by-hqztrue-s7wc/

3. 如果有一个数组的所有元素不相等，可以通过下标映射转换成 LIS， 值域在[0,max]的 严格递增 LIS 最优可以做到 O(nloglogmax)
