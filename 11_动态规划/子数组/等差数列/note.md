`446. 等差数列划分 II - 子序列 copy.py `
`1027. 最长等差数列.py`

`等差数列dp`：dp[i][diff] 表示以 nums[i] 结尾，公差为 diff 的等差数列的个数。

```Python
dp = [defaultdict(int) for _ in range(n)]
for i in range(1, n):
    for j in range(i):
        diff = nums[i] - nums[j]
        dp[i][diff] = dp[j][diff] + 1
return max(dp.values())
```

```Python
dp = defaultdict(lambda: 1)
for i in range(1, n):
    for j in range(i):
        diff = nums[i] - nums[j]
        dp[(i, diff)] = dp[(j, diff)] + 1
return max(dp.values())
```
