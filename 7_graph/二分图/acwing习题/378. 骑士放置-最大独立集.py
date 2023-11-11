# 给定一个 N×M 的棋盘，有一些格子禁止放棋子。
# !问棋盘上最多能放多少个不能互相攻击的骑士
# （国际象棋的“骑士”，类似于中国象棋的“马”，按照“日”字攻击，但没有中国象棋“别马腿”的规则）。

# 第一行包含三个整数 N,M,T，其中 T 表示禁止放置的格子的数量。
# 接下来 T 行每行包含两个整数 x 和 y，表示位于第 x 行第 y 列的格子禁止放置，行列数从 1 开始。
# 输入格式

# !最大独立集 = 总点数 - 最大匹配 (注意还要减去禁止放置的点)

from hungarian import Hungarian

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

DIR = [(1, 2), (1, -2), (-1, 2), (-1, -2), (2, 1), (2, -1), (-2, 1), (-2, -1)]
if __name__ == "__main__":
    ROW, COL, t = map(int, input().split())
    bad = set()
    for _ in range(t):
        x, y = map(int, input().split())
        x, y = x - 1, y - 1
        bad.add(x * COL + y)

    H = Hungarian(ROW * COL, ROW * COL)
    # build graph
    for r in range(ROW):
        for c in range(COL):
            cur = r * COL + c
            if ((r + c) & 1) or (cur in bad):
                continue
            for dr, dc in DIR:
                nr, nc = r + dr, c + dc
                next = nr * COL + nc
                if 0 <= nr < ROW and 0 <= nc < COL and next not in bad:
                    H.addEdge(cur, next)

    # !最大独立集 注意减去坏的格子
    print(ROW * COL - H.work() - t)
