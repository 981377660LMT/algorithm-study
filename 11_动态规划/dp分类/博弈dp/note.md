记忆化搜索容易理解

多以 dfs(index)为搜索函数
`能使对方输就是自己赢，不能使对方输就是自己输`
`自己赢=自己赢+对手不赢`

```Python
 # 自己赢=自己赢或对手不赢
if curSum + select >= target or not dfs(curSum + select, visited | (1 << select)):
    return True
```

```JS
if (自己赢 || !dp[i - j]) {
  dp[i] = true
}
```

除了**石子游戏 VI**，其他《石子游戏》思路基本都一样
