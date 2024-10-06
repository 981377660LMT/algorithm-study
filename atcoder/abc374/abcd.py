from itertools import permutations, zip_longest
from math import sqrt
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
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
            t = 0
            curX, curY = 0, 0
            dist1 = 0
            dist2 = 0
            for i in range(N):
                j = perm[i]
                a, b, c, d = A[j], B[j], C[j], D[j]
                if s >> j & 1:
                    a, b, c, d = c, d, a, b
                dist1 += calc(curX, curY, a, b)
                dist2 += calc(a, b, c, d)
                curX, curY = c, d
            t = dist1 / S + dist2 / T
            res = min(res, t)

    print(res)
