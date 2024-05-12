ecurSum 这种参数由 visited 唯一确定
**因此时间复杂度仍然是 O(2^n)的**

[473. 火柴拼正方形-状压 dp](473.%20%E7%81%AB%E6%9F%B4%E6%8B%BC%E6%AD%A3%E6%96%B9%E5%BD%A2-%E7%8A%B6%E5%8E%8Bdp.py)

```Python
@lru_cache(None)
def dfs(visited: int, curSum: int) -> bool:
    """注意这里curSum由visited唯一确定 因此复杂度是O(n*2^n)"""
    if visited == (1 << n)-1:
        return True
    for i in range(n):
        if (visited >> i) & 1:
            continue
        if curSum + nums[i] <= div:
            if dfs(visited | (1 << i), (curSum + nums[i]) % div):
                return True
    return False
```
