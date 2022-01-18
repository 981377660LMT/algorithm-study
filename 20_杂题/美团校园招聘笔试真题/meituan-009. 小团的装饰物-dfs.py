from functools import lru_cache
from sys import setrecursionlimit

setrecursionlimit(10000)
MOD = 998244353

limit, target = list(map(int, input().split()))
# limit, target = 4, 2


@lru_cache(None)
def dfs(index: int, curCost: int) -> int:
    if index == target:
        return 1
    res = 0
    for nextCost in range(curCost, limit + 1, curCost):
        res += dfs(index + 1, nextCost)
        res %= MOD
    return res


res = dfs(0, 1) % MOD
dfs.cache_clear()
print(res)
