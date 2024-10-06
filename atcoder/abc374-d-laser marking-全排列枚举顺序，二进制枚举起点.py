# D - Laser Marking
# https://atcoder.jp/contests/abc374/tasks/abc374_d
# 全排列枚举顺序，二进制枚举起点0/1.
# !O(n!*2^n)


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


from itertools import permutations
from math import sqrt


INF = int(4e18)


if __name__ == "__main__":
    N, S, T = map(int, input().split())
    A, B, C, D = [], [], [], []
    for _ in range(N):
        a, b, c, d = map(int, input().split())
        A.append(a)
        B.append(b)
        C.append(c)
        D.append(d)

    def calc(x1, y1, x2, y2):
        return sqrt((x1 - x2) ** 2 + (y1 - y2) ** 2)

    res = INF
    for perm in permutations(range(N)):
        for s in range(1 << N):
            curX, curY = 0, 0
            dist1, dist2 = 0, 0
            for i in range(N):
                j = perm[i]
                a, b, c, d = A[j], B[j], C[j], D[j]  # (a, b) -> (c, d)
                if s >> j & 1:
                    a, b, c, d = c, d, a, b
                dist1 += calc(curX, curY, a, b)
                dist2 += calc(a, b, c, d)
                curX, curY = c, d
            res = min(res, dist1 / S + dist2 / T)

    print(res)
