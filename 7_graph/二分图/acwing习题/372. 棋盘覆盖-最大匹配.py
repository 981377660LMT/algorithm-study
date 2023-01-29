# 给定一个 N 行 N 列(N<=100)的棋盘，已知某些格子禁止放置。
# 求最多能往棋盘上放多少块的长度为 2、宽度为 1 的骨牌，骨牌的边界与格线重合（骨牌占用两个格子），并且任意两张骨牌都不重叠。

# 第一行包含两个整数 N 和 t，其中 t 为禁止放置的格子的数量。
# 接下来 t 行每行包含两个整数 x 和 y，表示位于第 x 行第 y 列的格子禁止放置，行列数从 1 开始。

from hungarian import Hungarian

import sys

sys.setrecursionlimit(int(1e9))


DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]
n, bad = map(int, input().split())
badSet = set()
for _ in range(bad):
    x, y = map(int, input().split())
    x, y = x - 1, y - 1
    badSet.add(x * n + y)

H = Hungarian(n * n, n * n)
for r in range(n):
    for c in range(n):
        cur = r * n + c
        if (r + c) & 1 or cur in badSet:
            continue
        for dr, dc in DIR4:
            nr, nc = r + dr, c + dc
            next = nr * n + nc
            if 0 <= nr < n and 0 <= nc < n and next not in badSet:
                H.addEdge(cur, next)

print(H.work())
