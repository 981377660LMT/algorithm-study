# 矩阵快速幂
# MatPow(base,mod,cacheLevel): 带缓存的矩阵快速幂,适合多次查询;
# matpow(m1,exp,mod): 普通的矩阵快速幂;
# matpow2(m1,exp,mod): numpy的矩阵快速幂(非常快);
# mul(m1,m2,mod): 矩阵乘法.


import numpy as np


def matqpow2(base: "np.ndarray", exp: int, mod: int) -> "np.ndarray":
    """np矩阵快速幂"""
    res = np.eye(*base.shape, dtype=np.uint64)
    while exp:
        if exp & 1:
            res = (res @ base) % mod
        base = (base @ base) % mod
        exp >>= 1
    return res


#####################################################


from typing import List

M = List[List[int]]


class MatPow:
    __slots__ = ("_n", "_mod", "_base", "_cacheLevel", "_useCache", "_cache")

    def __init__(self, base: M, mod=1000000007, cacheLevel=4):
        n = len(base)
        b = [0] * n * n
        for i in range(n):
            for j in range(n):
                b[i * n + j] = base[i][j]
        useCache = cacheLevel >= 2
        self._n = n
        self._mod = mod
        self._base = b
        self._cacheLevel = cacheLevel
        self._useCache = useCache
        if useCache:
            self._cache = [[] for _ in range(cacheLevel - 1)]

    def pow(self, exp: int) -> M:
        if not self._useCache:
            return self._powWithOutCache(exp)
        if len(self._cache[0]) == 0:
            self._cache[0].append(self._base)
            for i in range(1, self._cacheLevel - 1):
                self._cache[i].append(self._mul(self._cache[i - 1][0], self._base))

        e = self._eye(self._n)
        div = 0
        while exp:
            if div == len(self._cache[0]):
                self._cache[0].append(
                    self._mul(self._cache[self._cacheLevel - 2][div - 1], self._cache[0][div - 1])
                )
                for i in range(1, self._cacheLevel - 1):
                    self._cache[i].append(self._mul(self._cache[i - 1][div], self._cache[0][div]))
            mod = exp % self._cacheLevel
            if mod:
                e = self._mul(e, self._cache[mod - 1][div])
            exp //= self._cacheLevel
            div += 1

        return self._to2D(e)

    def _mul(self, mat1: List[int], mat2: List[int]) -> List[int]:
        n = self._n
        res = [0] * n * n
        for i in range(n):
            for k in range(n):
                for j in range(n):
                    res[i * n + j] = (
                        res[i * n + j] + mat1[i * n + k] * mat2[k * n + j]
                    ) % self._mod
        return res

    def _powWithOutCache(self, exp: int) -> M:
        e = self._eye(self._n)
        b = self._base[:]
        while exp:
            if exp & 1:
                e = self._mul(e, b)
            exp >>= 1
            b = self._mul(b, b)
        return self._to2D(e)

    def _eye(self, n: int) -> List[int]:
        res = [0] * n * n
        for i in range(n):
            res[i * n + i] = 1
        return res

    def _to2D(self, mat: List[int]) -> M:
        n = self._n
        return [mat[i * n : (i + 1) * n] for i in range(n)]

    def __pow__(self, exp: int) -> M:
        return self.pow(exp)


def mul(mat1: M, mat2: M, mod: int) -> M:
    """矩阵相乘"""
    i_, j_, k_ = len(mat1), len(mat2[0]), len(mat2)
    res = [[0] * j_ for _ in range(i_)]
    for i in range(i_):
        for k in range(k_):
            for j in range(j_):
                res[i][j] = (res[i][j] + mat1[i][k] * mat2[k][j]) % mod
    return res


def matpow(base: M, exp: int, mod: int) -> M:
    e = [[0] * n for _ in range(n)]
    for i in range(n):
        e[i][i] = 1
    b = [row[:] for row in base]
    while exp:
        if exp & 1:
            e = mul(e, b, mod)
        exp >>= 1
        b = mul(b, b, mod)
    return e


if __name__ == "__main__":
    n = 876543210987654321
    MOD = int(1e9 + 7)

    dp = [[2], [1], [1]]  # 初始状态
    T = [[1, 1, 1], [1, 0, 0], [0, 1, 0]]
    mp = MatPow(T, MOD, cacheLevel=-1)
    resT = mp ** (n - 3)
    dp = mul(resT, dp, MOD)
    assert dp[0][0] == 639479200

    dp = [[2], [1], [1]]  # 初始状态
    T = np.array([[1, 1, 1], [1, 0, 0], [0, 1, 0]], np.uint64)
    resT = matqpow2(T, n - 3, MOD)
    dp = (resT @ dp) % MOD
    assert dp[0][0] == 639479200
