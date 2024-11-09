from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 num 。如果一个数字字符串的奇数位下标的数字之和与偶数位下标的数字之和相等，那么我们称这个数字字符串是 平衡的 。

# 请你返回 num 不同排列 中，平衡 字符串的数目。

# 由于答案可能很大，请你将答案对 109 + 7 取余 后返回。


# 一个字符串的 排列 指的是将字符串中的字符打乱顺序后连接得到的字符串。


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


E = Enumeration(100, MOD)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def countBalancedPermutations(self, num: str) -> int:
        n = len(num)
        counter = Counter(num)
        evenPos = (n + 1) // 2
        oddPos = n // 2
        totalSum = sum(int(c) for c in num)
        if totalSum % 2 != 0:
            return 0

        target = totalSum // 2
        dp = [defaultdict(int) for _ in range(evenPos + 1)]
        dp[0][0] = 1
        for digit in sorted(counter):
            d, c = int(digit), counter[digit]
            for i in range(evenPos, -1, -1):
                for s in list(dp[i].keys()):
                    for k in range(1, min2(c, evenPos - i) + 1):
                        newI, newS = i + k, s + k * d
                        if newS > target:
                            break
                        dp[newI][newS] = (dp[newI][newS] + dp[i][s] * E.C(c, k)) % MOD
        ways = dp[evenPos][target]
        if ways == 0:
            return 0
        deno = 1
        for c in counter.values():
            deno = deno * E.ifac(c) % MOD
        return ways * E.fac(evenPos) * E.fac(oddPos) * deno % MOD