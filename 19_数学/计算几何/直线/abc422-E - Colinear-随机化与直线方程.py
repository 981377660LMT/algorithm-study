# E - Colinear-随机化、直线方程
# https://atcoder.jp/contests/abc422/tasks/abc422_e
#
# 在二维平面上有 N 个点。N 是一个奇数。第 i 个点位于 (x_i, y_i)。所有点的坐标都是不同的。
#
# 请你判断是否存在一条直线，它能穿过这 N 个点中的过半数（即至少 (N+1)/2 个点）。如果存在，请输出这条直线。
#
# 另外，对于满足限制条件的任何输入，如果存在满足条件的直线，那么这条直线可以用 ax + by + c = 0 的形式表示，其中 a, b, c 是 -10^18 到 10^18 之间的整数（且 (a,b,c) ≠ (0,0,0)）。请你输出这样的 a, b, c。
#
# 限制条件
#
# 3 ≤ N ≤ 5 * 10^5
# N 是奇数
# -10^8 ≤ x_i ≤ 10^8
# -10^8 ≤ y_i ≤ 10^8
# 如果 i ≠ j，则 (x_i, y_i) ≠ (x_j, y_j)
# 所有输入值均为整数
#
# !如果存在一条直线穿过超过半数的点（我们称之为“多数派直线”），那么我们随机选择两个不同的点，这两个点都属于这条多数派直线的概率非常高。

import random
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N = int(input())
    X, Y = [0] * N, [0] * N
    for i in range(N):
        X[i], Y[i] = map(int, input().split())

    iterations = min(100, N * (N - 1) // 2)
    for _ in range(iterations):
        i, j = random.sample(range(N), 2)

        x1, y1 = X[i], Y[i]
        x2, y2 = X[j], Y[j]

        # 直线方程系数
        a = y2 - y1
        b = x1 - x2
        c = -a * x1 - b * y1

        count = 0
        for x, y in zip(X, Y):
            if a * x + b * y + c == 0:
                count += 1

        if count > N // 2:
            print("Yes")
            print(a, b, c)
            exit(0)

    print("No")
