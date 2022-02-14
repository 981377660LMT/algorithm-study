# 给定三种类型的小球 P、Q、R，每种小球的数量分别为 np、nq、nr 个
# 现在想将这些小球排成一条直线，但是`不允许相同类型的小球相邻`，问有多少种排列方法
from functools import lru_cache


@lru_cache(None)
def dfs(pre: int, np: int, nq: int, nr: int) -> int:
    if np + nq + nr == 0:
        return 1

    res = 0
    if pre != 0 and np > 0:
        res += dfs(0, np - 1, nq, nr)
    if pre != 1 and nq > 0:
        res += dfs(1, np, nq - 1, nr)
    if pre != 2 and nr > 0:
        res += dfs(2, np, nq, nr - 1)

    return res


np, nq, nr = [int(x) for x in input().split()]
res = 0
if np > 0:
    res += dfs(0, np - 1, nq, nr)
if nq > 0:
    res += dfs(1, np, nq - 1, nr)
if nr > 0:
    res += dfs(2, np, nq, nr - 1)
print(res)

