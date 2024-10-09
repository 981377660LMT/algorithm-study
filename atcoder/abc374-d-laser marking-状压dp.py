# D - Laser Marking
# https://atcoder.jp/contests/abc374/tasks/abc374_d
#
# !O(n^2*2^n)


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


from functools import lru_cache
from math import sqrt


INF = int(4e18)


def min2(a, b):
    return a if a < b else b


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

    T1 = [calc(A[i], B[i], C[i], D[i]) / T for i in range(N)]

    MASK = (1 << N) - 1

    @lru_cache(None)
    def dfs(visited: int, x: int, y: int) -> float:
        if visited == MASK:
            return 0
        res = INF
        for next_ in range(N):
            if (visited >> next_) & 1:
                continue
            a, b, c, d = A[next_], B[next_], C[next_], D[next_]
            res = min2(res, calc(x, y, a, b) / S + T1[next_] + dfs(visited | (1 << next_), c, d))
            res = min2(res, calc(x, y, c, d) / S + T1[next_] + dfs(visited | (1 << next_), a, b))
        return res

    res = dfs(0, 0, 0)
    dfs.cache_clear()
    print(res)
