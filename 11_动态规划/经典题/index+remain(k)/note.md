**k/remain 维度 一般是指的在某个序列中进行某种操作 k/remain 次**

`注意 index == n 要返回`

```Python
@lru_cache(None)
def dfs(index: int, remain: int) -> int:
    if index == n or remain < 0:
      return 0

    res = 0
    for select in range(min(remain + 1, len(piles[index]) + 1)):  # !这里可以根据题目优化
        next = dfs(index + 1, remain - select)
        res = max(res, next + preSums[index][select])
    return res
```

**注意这种题如果是求和，可以用前缀(和)优化 dp 范围转移的复杂度**
[!注意到最内层的转移可以前缀和优化 => dp 由一连串的 index 转移过来 所以考虑把 index 作为第二维度遍历](./6244.%20完美分割的方案数-前缀和优化.ts)
