# 常系数线性递推
# https://tjkendev.github.io/procon-library/python/series/kitamasa.html
# !O(k^2logn) 求线性递推式的第n项 (比矩阵快速幂快一个k)
# 線形漸化式 dp[i+k] = c0*dp[i] + c1*dp[i+1] + ... + ci+k-1*dp[i+k-1] (i>=0) の第n項を求める
# C: 系数 c0,c1,...,ci+k-1
# A: dp[0]-dp[k-1] 初始值
# n: 第n项

from typing import List


MOD = int(1e9 + 7)


def kitamasa(C: List[int], A: List[int], n: int) -> int:
    if n == 0:
        return A[0]

    assert len(C) == len(A)
    k = len(C)
    C0 = [0] * k
    C1 = [0] * k
    C0[1] = 1

    def inc(k, C0, C1):
        C1[0] = C0[k - 1] * C[0] % MOD
        for i in range(k - 1):
            C1[i + 1] = (C0[i] + C0[k - 1] * C[i + 1]) % MOD

    def dbl(k, C0, C1):
        D0 = [0] * k
        D1 = [0] * k
        D0[:] = C0[:]
        for j in range(k):
            C1[j] = C0[0] * C0[j] % MOD
        for i in range(1, k):
            inc(k, D0, D1)
            for j in range(k):
                C1[j] += C0[i] * D1[j] % MOD
            D0, D1 = D1, D0
        for i in range(k):
            C1[i] %= MOD

    p = n.bit_length() - 1

    while p:
        p -= 1
        dbl(k, C0, C1)
        C0, C1 = C1, C0
        if (n >> p) & 1:
            inc(k, C0, C1)
            C0, C1 = C1, C0

    res = 0
    for i in range(k):
        res = (res + C0[i] * A[i]) % MOD
    return res


# 斐波那契
def fib(n: int) -> int:
    """0 1 1 2 3 5 8 13 21 34 55"""
    return kitamasa([1, 1], [0, 1], n)


K, N = map(int, input().split())
print(kitamasa([1] * K, [1] * K, N - 1))
