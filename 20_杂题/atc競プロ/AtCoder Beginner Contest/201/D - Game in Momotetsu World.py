"""
桃太郎电铁游戏里有+-两种格子
开始时左上角有一个棋子，每次可以向右或向下移动一格
移动到+格子得1分,移动到-格子得-1分
每个人都最佳移动策略,问最后谁会赢 

ROW,COL<=2000 博弈dp计算分数差值
"""
from functools import lru_cache
from typing import List

INF = int(1e9)


def gameInMomotetsuWorld(grid: List[str]) -> int:
    @lru_cache(None)
    def dfs(row: int, col: int) -> int:
        if row == ROW - 1 and col == COL - 1:
            return 0
        res = -INF
        if row + 1 < ROW:
            res = max(res, (1 if grid[row + 1][col] == "+" else -1) - dfs(row + 1, col))
        if col + 1 < COL:
            res = max(res, (1 if grid[row][col + 1] == "+" else -1) - dfs(row, col + 1))
        return res

    ROW, COL = len(grid), len(grid[0])
    res = dfs(0, 0)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    ROW, COL = map(int, input().split())
    grid = [input() for _ in range(ROW)]
    res = gameInMomotetsuWorld(grid)
    if res > 0:
        print("Takahashi")
    elif res < 0:
        print("Aoki")
    else:
        print("Draw")
