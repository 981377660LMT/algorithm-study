`一种理解：dp 是对集合的划分，在各个集合中求最优解`

dp[i][n]数组表示现在的**现金最大值**: i 代表天数，n 代表是股票状态

2\*k+1 个状态
没有操作
第一次买入
第一次卖出
第二次买入
第二次卖出
...
**除了 0 以外，偶数就是卖出，奇数就是买入。**

```JS
const dp = Array.from({ length: len }, () => Array(k).fill(0))
dp[0][0] = 0
dp[0][1] = -prices[0]
dp[0][2] = -Infinity
dp[0][3] = -Infinity
dp[0][4] = -Infinity
...

for (let i = 1; i < len; i++) {
  for (let j = 0; j < 2 * k; j += 2) {
    dp[i][j + 1] = Math.max(dp[i - 1][j + 1], dp[i - 1][j] - prices[i])
    dp[i][j + 2] = Math.max(dp[i - 1][j + 2], dp[i - 1][j + 1] + prices[i])
  }
}


```

[总结](https://programmercarl.com/%E5%8A%A8%E6%80%81%E8%A7%84%E5%88%92-%E8%82%A1%E7%A5%A8%E9%97%AE%E9%A2%98%E6%80%BB%E7%BB%93%E7%AF%87.html#%E5%8D%96%E8%82%A1%E7%A5%A8%E7%9A%84%E6%9C%80%E4%BD%B3%E6%97%B6%E6%9C%BA)
