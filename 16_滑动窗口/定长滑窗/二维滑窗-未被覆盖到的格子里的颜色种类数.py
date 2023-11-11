# !二维滑窗-未被覆盖到的格子里的颜色种类数
# !没有被覆盖到的格子里的颜色数

from collections import defaultdict
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(
    ROW: int, COL: int, windowRow: int, windowCol: int, grid: List[List[int]]
) -> List[List[int]]:
    """返回每个windowRow*winowCol窗口之外的格子里的颜色种类数
    ROW,COL <= 300
    """

    def cal(row: int) -> List[int]:
        # !二维滑动窗口
        res, fullCount = [], 0
        windowCounter = defaultdict(int)
        for right in range(COL):
            for r in range(row, windowRow + row):
                windowCounter[grid[r][right]] += 1
                if windowCounter[grid[r][right]] == gridCounter[grid[r][right]]:
                    fullCount += 1
            if right >= windowCol:
                for r in range(row, windowRow + row):
                    windowCounter[grid[r][right - windowCol]] -= 1
                    if (
                        windowCounter[grid[r][right - windowCol]]
                        == gridCounter[grid[r][right - windowCol]] - 1
                    ):
                        fullCount -= 1
            if right >= windowCol - 1:
                res.append(gridCount - fullCount)
        return res

    gridCounter = defaultdict(int)
    for i in range(ROW):
        for j in range(COL):
            gridCounter[grid[i][j]] += 1
    gridCount = len(gridCounter)

    return [cal(r) for r in range(ROW - windowRow + 1)]


if __name__ == "__main__":
    ROW, COL, n, windowRow, windowCol = map(int, input().split())
    grid = []
    for i in range(ROW):
        grid.append(list(map(int, input().split())))
    res = solve(ROW, COL, windowRow, windowCol, grid)
    for row in res:
        print(*row)
