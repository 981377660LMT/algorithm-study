"""
https://atcoder.jp/contests/abc231/submissions/27959643
棋盘上有n个白色的马
支付costi的费用,可以将第i个马(xi,yi)变成黑色的马
要使每一行每一列都至少有一个黑色的马
问最少需要支付多少费用(题目数据保证有解)
ROW,COL<=1e3
n<=1e3

带权的最小边覆盖
"""

import sys
from MinCostMaxFlow import MinCostMaxFlowDinic

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    ROW, COL, n = map(int, input().split())
    horse = []
    for _ in range(n):
        r, c, cost = map(int, input().split())
        r, c = r - 1, c - 1
        horse.append((r, c, cost))

    START, END = ROW + COL + 4, ROW + COL + 5
    # mcmf = MinCostMaxFlow(ROW + COL + 10, START, END)
    # for i in range(ROW):
    #     mcmf.addEdge(START, i, rowDeg[i] - 1, 0)
    # for i in range(COL):
    #     mcmf.addEdge(ROW + i, END, colDeg[i] - 1, 0)
    # for r, c, cost in horse:
    #     mcmf.addEdge(r, ROW + c, 1, cost)
    # print(mcmf.work()[1])


# TODO
