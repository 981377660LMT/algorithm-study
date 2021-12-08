`446. 等差数列划分 II - 子序列 copy.py `
`1027. 最长等差数列.py`

```Python
dp = defaultdict(lambda: 1)
for i in range(1, n):
    for j in range(i):
        diff = nums[i] - nums[j]
        dp[(i, diff)] = dp[(j, diff)] + 1
return max(dp.values())
```
