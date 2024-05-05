from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(998244353)
INF = int(1e20)

# 一个二进制字符串是由 0 和 1 组成的字符串。如果一个长度为 m（m 为奇数）的二进制字符串 s 满足 最后一位数字 与 s 的「排序中间数」相同，那么它就是美丽的。如 11001 是美丽的，而 11100 不是美丽的。


# 给定一个长度为 n 的二进制字符串 s，请返回 s 所有前缀的美丽扩展的个数之和。


# 由于答案可能很大，你只需要求出它模 998244353 的结果。
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


C = Enumeration(int(2e5 + 10), MOD)


class Solution:
    def beautifulString(self, s: str) -> int:
        curOne, curZero = 0, 0
        res = 0
        for i, c in enumerate(s):
            if c == "1":
                curOne += 1
            else:
                curZero += 1

            # 从i个位置选择若干个数，使得c的个数严格>=i+1
            atLeastSelect = (i + 1) - (curOne if c == "1" else curZero)
            # [0,atLeastSelect-1]
            res += pow(2, i, MOD) - sum(C.C(i, j) for j in range(atLeastSelect))
            res %= MOD

        return res % MOD
