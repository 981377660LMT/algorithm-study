# 活字印刷
# 你有一套活字字模 tiles，其中每个字模上都刻有一个字母 tiles[i]。
# 返回你可以印出的非空字母序列的数目。
# 注意：本题中，每个活字字模只能使用一次。
# 1 <= tiles.length <= 7

from collections import Counter
from itertools import permutations

MOD = int(1e9 + 7)


class Enumeration:
    __slots__ = ("_fac", "_ifac", "_inv", "_mod")

    def __init__(self, size: int, mod: int) -> None:
        self._mod = mod
        self._fac = [1]
        self._ifac = [1]
        self._inv = [1]
        self._expand(size)

    def fac(self, k: int) -> int:
        self._expand(k)
        return self._fac[k]

    def ifac(self, k: int) -> int:
        self._expand(k)
        return self._ifac[k]

    def inv(self, k: int) -> int:
        """模逆元"""
        self._expand(k)
        return self._inv[k]

    def C(self, n: int, k: int) -> int:
        if n < 0 or k < 0 or n < k:
            return 0
        mod = self._mod
        return self.fac(n) * self.ifac(k) % mod * self.ifac(n - k) % mod

    def P(self, n: int, k: int) -> int:
        if n < 0 or k < 0 or n < k:
            return 0
        mod = self._mod
        return self.fac(n) * self.ifac(n - k) % mod

    def H(self, n: int, k: int) -> int:
        """可重复选取元素的组合数"""
        if n == 0:
            return 1 if k == 0 else 0
        return self.C(n + k - 1, k)

    def put(self, n: int, k: int) -> int:
        """n个相同的球放入k个不同的盒子(盒子可放任意个球)的方案数."""
        return self.C(n + k - 1, n)

    def catalan(self, n: int) -> int:
        """卡特兰数"""
        return self.C(2 * n, n) * self.inv(n + 1) % self._mod

    def _expand(self, size: int) -> None:
        size = min(size, self._mod - 1)
        if len(self._fac) < size + 1:
            mod = self._mod
            preSize = len(self._fac)
            diff = size + 1 - preSize
            self._fac += [1] * diff
            self._ifac += [1] * diff
            self._inv += [1] * diff
            for i in range(preSize, size + 1):
                self._fac[i] = self._fac[i - 1] * i % mod
            self._ifac[size] = pow(self._fac[size], mod - 2, mod)  # !modInv
            for i in range(size - 1, preSize - 1, -1):
                self._ifac[i] = self._ifac[i + 1] * (i + 1) % mod
            for i in range(preSize, size + 1):
                self._inv[i] = self._ifac[i] * self._fac[i - 1] % mod


E = Enumeration(int(1e3), MOD)


class Solution:
    def numTilePossibilities(self, tiles: str) -> int:
        """
        O(n^2) dp.
        dp[i][j]表示用前i种字母组成长度为j的序列的个数.
        如果选第i种字母k个, 则dp[i][j] = dp[i-1][j-k]*C(j, k).
        """
        freq = Counter(tiles).values()  # 每个字母的频率
        dp = [[0] * (len(tiles) + 1) for _ in range(len(freq) + 1)]
        dp[0][0] = 1
        for i, count in enumerate(freq, 1):
            for j in range(len(tiles) + 1):
                for k in range(min(j, count) + 1):
                    dp[i][j] += dp[i - 1][j - k] * E.C(j, k)
                    dp[i][j] %= MOD
        return sum(dp[-1][1:]) % MOD

    def numTilePossibilities2(self, tiles: str) -> int:
        """O(n!*n)."""
        res = set()
        for len_ in range(1, len(tiles) + 1):
            for perm in permutations(tiles, len_):
                res.add(perm)
        return len(res)


print(Solution().numTilePossibilities("AAB"))
# 输出：8
# 解释：可能的序列为 "A", "B", "AA", "AB", "BA", "AAB", "ABA", "BAA"。
