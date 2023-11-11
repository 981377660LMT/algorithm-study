# 蒙德里安的梦想
# !求把 N×M 的棋盘分割成若干个 1×2 的长方形，有多少种方案。
# n,m<=11

# 轮廓线dp
# !技巧:多加一个虚拟行 用来检查最后一行的状态 即有没有放满

from functools import lru_cache
import sys
from typing import Tuple

sys.setrecursionlimit(int(1e6))


def solve(row: int, col: int) -> int:
    """轮廓线dp 时间复杂度O(n*m*2*min(n,m))"""

    @lru_cache(None)
    def dfs(r: int, c: int, cState: Tuple[int, ...]) -> int:
        if r == ROW:  # !检查最后一行是否每个格子被占满
            return 1 if all(cState) else 0
        if c == COL:
            return dfs(r + 1, 0, cState)

        res = 0

        # 三种转移:

        # 1.放竖的
        if r > 0 and cState[0] == 0:
            res += dfs(r, c + 1, cState[1:] + (1,))

        # 2.放横的 (前提是上面格子被放了 且 左边格子没被放)
        if (r == 0 or cState[0] == 1) and (c > 0 and cState[-1] == 0):
            res += dfs(r, c + 1, cState[1:-1] + (1, 1))

        # 3.不放 (前提是上面格子被放了)
        if r == 0 or cState[0] == 1:
            res += dfs(r, c + 1, cState[1:] + (0,))

        return res

    ROW, COL = sorted([row, col], reverse=True)
    # ROW += 1  # !多一行用来检验最后一行是否每个格子都被放了
    res = dfs(0, 0, tuple([1] * COL))  # 一开始放不了
    dfs.cache_clear()
    return res


while True:
    n, m = map(int, input().split())
    if 0 in (n, m):
        break
    print(solve(n, m))
