`注意 index == n 要返回啊!`

```Python
@lru_cache(None)
def dfs(index: int, remain: int) -> int:
    if index == n or remain < 0:
      return 0

    res = 0
    for select in range(min(remain + 1, len(piles[index]) + 1)):
        next = dfs(index + 1, remain - select)
        res = max(res, next + preSums[index][select])
    return res
```
