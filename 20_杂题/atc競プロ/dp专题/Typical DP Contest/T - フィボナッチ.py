# 求k斐波那契数列的第n项
# k<=1000
# n<=1e9


import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def kbonacci2(k: int, n: int) -> int:
    """k-bonacci number 第n项

    时间复杂度: O(k^2k*logn)

    >>> a1=a2=...=aK=1
    >>> ai=sum(aj for j in range(i-k,i))

    长为n的01序列中1不能连续出现k次的方案数
    """
    assert k >= 1, n >= 0

    def linear_recursion_solver(a: List[int], x: List[int], k: int, e0: int, e1: int) -> int:
        """https://atcoder.jp/contests/tdpc/submissions/15359686"""

        def rec(k: int) -> List[int]:
            c = [e0] * m
            if k < m:
                c[k] = e1
                return c[:]
            b = rec(k // 2)
            t = [e0] * (2 * m + 1)
            for i in range(m):
                for j in range(m):
                    t[i + j + (k & 1)] += b[i] * b[j]
                    t[i + j + (k & 1)] %= MOD
            for i in reversed(range(m, 2 * m)):
                for j in range(m):
                    t[i - m + j] += a[j] * t[i]
                    t[i - m + j] %= MOD
            for i in range(m):
                c[i] = t[i]
            return c[:]

        m = len(a)
        c = rec(k)
        res = 0
        for ci, xi in zip(c, x):
            res += ci * xi
            res %= MOD
        return res

    A, C = [1] * k, [1] * k
    return linear_recursion_solver(C[::-1], A, n, 0, 1) % MOD


k, n = map(int, input().split())
print(kbonacci2(k, n - 1))
