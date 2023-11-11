# 在好的格子('.')上放置马 使得马之间不互相攻击(不出现在与自己相邻的8个方向) 求摆放的方案数
# '.'表示空地 '#'表示障碍
# ROW,COL<=24

# 轮廓线dp的优化
#  `不存在连续两个元素相同`的状态 由此可以剪枝
# `将状态数从 2^COL 变为 fib(COL) 约等于1.62^COL 大幅度剪枝`
from typing import Tuple
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


@lru_cache(None)
def dfs(r: int, c: int, cState: Tuple[int, ...]) -> int:
    """时间复杂度O(n * m * (1.62 ^ m))"""
    if r == ROW:
        return 1
    if c == COL:
        return dfs(r + 1, 0, cState)

    res = dfs(r, c + 1, cState[1:] + (0,))  # 不选择当前位置
    if matrix[r][c] == '.':
        """不能与相邻八个方向的马紧挨着"""
        leftUp = cState[0] if r and c else 0
        upUp = cState[1] if r else 0
        rightUp = cState[2] if r and c + 1 < COL else 0
        left = cState[-1] if c else 0
        if leftUp == upUp == rightUp == left == 0:
            res += dfs(r, c + 1, cState[1:] + (1,))
            res %= MOD

    return res


ROW, COL = map(int, input().split())
matrix = []
for _ in range(ROW):
    matrix.append(list(input()))  # '.'表示空地 '#'表示障碍

res = dfs(0, 0, tuple([0] * (COL + 1)))
dfs.cache_clear()
print(res)
