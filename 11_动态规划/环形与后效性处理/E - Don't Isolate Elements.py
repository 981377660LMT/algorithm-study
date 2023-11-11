from functools import lru_cache
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 可以反转任意的行使得所有的格子都不是孤立的格子(即与上下左右相同的格子不同)
# 求最少的反转次数
# !dp[i][j][k]:考虑到第i行，第i-1行翻/没翻，第i行翻/没翻的前i-1行均合法的最小操作次数
# !如果dp有后效性(需要下一行的状态),那么可以枚举下一行的状态


def dontIsolateElements(grid: List[List[int]]) -> int:
    @lru_cache(None)
    def dfs(row: int, preFlip: int, curFlip: int) -> int:
        """前一行是否翻转，当前行是否翻转"""
        if row == ROW:
            return 0
        res = INF
        for nextFlip in range(2):
            ok1 = True
            for col in range(COL):
                ok2 = False
                if col > 0:
                    ok2 |= (grid[row][col] ^ curFlip) == (grid[row][col - 1] ^ curFlip)
                if col < COL - 1:
                    ok2 |= (grid[row][col] ^ curFlip) == (grid[row][col + 1] ^ curFlip)
                if row > 0:
                    ok2 |= (grid[row][col] ^ curFlip) == (grid[row - 1][col] ^ preFlip)
                if row < ROW - 1:
                    ok2 |= (grid[row][col] ^ curFlip) == (grid[row + 1][col] ^ nextFlip)
                if not ok2:
                    ok1 = False
                    break
            if ok1:
                res = min(res, dfs(row + 1, curFlip, nextFlip) + nextFlip)
        return res

    ROW, COL = len(grid), len(grid[0])
    res = min(dfs(0, 0, 0), dfs(0, 0, 1) + 1)
    dfs.cache_clear()
    return res if res < INF else -1


if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(ROW)]
    print(dontIsolateElements(grid))
