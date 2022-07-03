# 前缀和优化 dp dp[i]=dp[left]...+dp[i-1]

- 一般会用到 dp 和 dpSum 数组
- 求出 dp[i]后再更新 dpSum[i]
  [分割数组使得逆序数不大于 k](089%20-%20Partitions%20and%20Inversions%EF%BC%88%E2%98%857%EF%BC%89.py)
  [1997. 访问完所有房间的第一天-前缀和优化](1997.%20%E8%AE%BF%E9%97%AE%E5%AE%8C%E6%89%80%E6%9C%89%E6%88%BF%E9%97%B4%E7%9A%84%E7%AC%AC%E4%B8%80%E5%A4%A9-%E5%89%8D%E7%BC%80%E5%92%8C%E4%BC%98%E5%8C%96.py)

```Python
# dp[i]=dp[lefti]...+dp[i-1] 转化成
dp[i] = dpSum[i - 1] - dpSum[lefti - 1]
dpSum[i] = dpSum[i - 1] + dp[i]
```

## 前缀和 dp 的题目特点是某个点的增量/总贡献来源于之前的某段区间
