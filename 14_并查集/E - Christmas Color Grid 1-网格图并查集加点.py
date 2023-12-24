# https://atcoder.jp/contests/abc334/tasks/abc334_e
# E - Christmas Color Grid 1-网格图并查集加点(最大人工岛)
# 给定一个01矩阵.
# !对每个0，将其变为1后(加点)，求图中1组成的联通分量个数.


from typing import List
from UnionFindWithUndo import UnionFindArrayWithUndo

DIR4 = ((-1, 0), (1, 0), (0, -1), (0, 1))


def christmasColorGrid1(grid: List[List[int]]) -> List[int]:
    ROW, COL = len(grid), len(grid[0])
    uf = UnionFindArrayWithUndo(ROW * COL)

    onesCount = 0
    for x in range(ROW):
        for y in range(COL):
            if grid[x][y] == 0:
                continue
            onesCount += 1
            cur = x * COL + y
            if x > 0 and grid[x - 1][y] == 1:
                uf.union(cur, (x - 1) * COL + y)
            if y > 0 and grid[x][y - 1] == 1:
                uf.union(cur, x * COL + y - 1)
    zeroCount = ROW * COL - onesCount

    res = []
    state = uf.getState()
    for x in range(ROW):
        for y in range(COL):
            if grid[x][y] == 1:
                continue
            cur = x * COL + y
            for dx, dy in DIR4:
                nc, ny = x + dx, y + dy
                if 0 <= nc < ROW and 0 <= ny < COL and grid[nc][ny] == 1:
                    uf.union(cur, nc * COL + ny)
            res.append(uf.part - zeroCount + 1)
            uf.rollback(state)
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    MOD = 998244353

    h, w = map(int, input().split())
    grid = []
    zeroCount = 0
    for _ in range(h):
        s = input()
        zeroCount += s.count(".")
        grid.append([0 if c == "." else 1 for c in s])

    res = christmasColorGrid1(grid)
    print(sum(res) * pow(zeroCount, MOD - 2, MOD) % MOD)
