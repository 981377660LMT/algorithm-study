# 1x2铺瓷砖输出最大方案
# https://atcoder.jp/contests/practice2/tasks/practice2_d
# 给定一个二维网格,"."表示空地,"#"表示障碍物,要求用1x2的瓷砖尽可能多地铺满空地.
# 输出最大铺瓷砖方案数,以及任意一种铺瓷砖方案.
# !上下瓷砖用'v'和'^'表示,左右瓷砖用'>'和'<'表示.
# ROW,COL<=100

from 匈牙利算法 import Hungarian

from typing import List, Tuple


DIR4 = ((-1, 0), (1, 0), (0, -1), (0, 1))
UP = "^"
DOWN = "v"
LEFT = "<"
RIGHT = ">"


def solve(grid: List[List[str]]) -> Tuple[int, List[List[str]]]:
    ROW, COL = len(grid), len(grid[0])
    flow = Hungarian()
    for r in range(ROW):
        for c in range(COL):
            if (r + c) & 1 or grid[r][c] == "#":  # left: 偶数格子 right: 奇数格子
                continue
            for dr, dc in DIR4:
                nr, nc = r + dr, c + dc
                if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == ".":
                    flow.addEdge(r * COL + c, nr * COL + nc)

    res = flow.work()
    solution = [list(s) for s in grid]
    for left, right in res:
        r1, c1 = left // COL, left % COL
        r2, c2 = right // COL, right % COL
        if r1 + 1 == r2:
            solution[r1][c1], solution[r2][c2] = DOWN, UP
        elif r1 - 1 == r2:
            solution[r1][c1], solution[r2][c2] = UP, DOWN
        elif c1 + 1 == c2:
            solution[r1][c1], solution[r2][c2] = RIGHT, LEFT
        elif c1 - 1 == c2:
            solution[r1][c1], solution[r2][c2] = LEFT, RIGHT

    return len(res), solution


if __name__ == "__main__":
    n, m = map(int, input().split())
    grid = [list(input()) for _ in range(n)]
    res, solution = solve(grid)
    print(res)
    for row in solution:
        print("".join(row))

    # 输入:
    # 3 3
    # #..
    # ..#
    # ...

    # 输出:
    # 3
    # #><
    # vv#
    # ^^.
