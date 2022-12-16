# 斑马条纹(黑白格子相邻不同色的对数的最大值)
# 给定n × n矩阵，每个格子为B,W或者?，
# ?可以填任意颜色。求相邻格子不同色的对数的最大值。
# n<=100

# !先假设所有相邻边都不同色 变成最小割问题
# 1.最少代价:计算扣分 注意需要奇偶翻转判断
# 2.所有格子分成两组(白/黑)
# 源点:白色格子((r+c)&1==0) 黑色格子((r+c)&1==1)
# 汇点:黑色格子((r+c)&1==0) 白色格子((r+c)&1==1)
# !从而每个点需要与其相邻的点连边 代价为1(如果同色 就必须割去这条边 即少1个相邻不同色的对数)
# 3.不合法的方案
# !原来的白色格子必须是白色(START) 原来的黑色格子必须是黑色(END)

# https://kanpurin.hatenablog.com/entry/2021/02/27/225330

# TODO
from typing import List
from MaxFlow import ATCMaxFlow


DIR4 = ((-1, 0), (0, 1), (1, 0), (0, -1))
INF = int(1e18)


def zebraness(grid: List[str]) -> int:
    n = len(grid)
    START, END = n * n, n * n + 1
    maxFlow = ATCMaxFlow(END + 1, START, END)

    for r in range(n):
        for c in range(n):
            cur, pos = grid[r][c], r * n + c
            for dr, dc in DIR4:
                nr, nc = r + dr, c + dc
                if 0 <= nr < n and 0 <= nc < n:
                    maxFlow.addEdge(pos, nr * n + nc, 1)  # !1表示如果同色 就必须割去这条边 即少1个相邻不同色的对数
            if cur == "?":
                continue
            if (cur == "W") ^ ((r + c) & 1):
                maxFlow.addEdge(START, pos, INF)  # !原来白色格子必须是白色(START)
            else:
                maxFlow.addEdge(pos, END, INF)  # !原来黑色格子必须是黑色(END)

    return n * (n - 1) * 2 - maxFlow.calMaxFlow()


n = int(input())
grid = [input() for _ in range(n)]
print(zebraness(grid))
