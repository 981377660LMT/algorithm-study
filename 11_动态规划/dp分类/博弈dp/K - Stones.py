# n堆石子 石子总数为k
# !每次选择一个石子数 从堆里移除这么多个石子
# 不能移除的人输

# k<=1e5 n<=100
# 博弈dp
from functools import lru_cache
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, k = map(int, input().split())
nums = list(map(int, input().split()))
WIN1, WIN2 = "First", "Second"


@lru_cache(None)
def dfs(remain: int) -> bool:
    if remain == 0:
        return False
    res = False
    for num in nums:
        if remain >= num:
            res |= not dfs(remain - num)
    return res


res = dfs(k)
dfs.cache_clear()
print(WIN1 if res else WIN2)
