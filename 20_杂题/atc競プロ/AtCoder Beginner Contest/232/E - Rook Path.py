# 机器人每一步可以移动到同一行，或者同一列上的任意一点（除了自己现在呆的点）。
# 问从点(x1, y1)到点(x2,y2)用了k步有几种方案。
# !dp预处理
# ROW,COL<=1e9
# k<=1e6

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    ROW, COL, k = map(int, input().split())
    x1, y1, x2, y2 = map(int, input().split())
