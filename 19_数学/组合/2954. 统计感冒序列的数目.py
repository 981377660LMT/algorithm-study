# 2954. 统计感冒序列的数目
# https://leetcode.cn/problems/count-the-number-of-infection-sequences/solutions/2551734/zu-he-shu-xue-ti-by-endlesscheng-5fjp/
# 给你一个整数 n 和一个下标从 0 开始的整数数组 sick ，数组按 升序 排序。
# 有 n 位小朋友站成一排，按顺序编号为 0 到 n - 1 。
# 数组 sick 包含一开始得了感冒的小朋友的位置。
# 如果位置为 i 的小朋友得了感冒，他会传染给下标为 i - 1 或者 i + 1 的小朋友，前提 是被传染的小朋友存在且还没有得感冒。每一秒中， 至多一位 还没感冒的小朋友会被传染。
# 经过有限的秒数后，队列中所有小朋友都会感冒。感冒序列 指的是 所有 一开始没有感冒的小朋友最后得感冒的顺序序列。请你返回所有感冒序列的数目。
# 由于答案可能很大，请你将答案对 109 + 7 取余后返回。
# 注意，感冒序列 不 包含一开始就得了感冒的小朋友的下标
#
#
# !每段空白区间无关，若不在两端答案为 $2 ^ {len - 1}$
# 段与段之间的关系是独立的，把所有放法相乘，再乘上每种感冒序列的方案，即为答案

from typing import List


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


MOD = int(1e9 + 7)
E = Enumeration(int(1e5) + 10, MOD)


class Solution:
    def numberOfSequence(self, n: int, sick: List[int]) -> int:
        res = 1
        m = n - len(sick)
        for x, y in zip(sick, sick[1:]):
            v = y - x - 1
            if v:
                res = res * E.C(m, v) * pow(2, v - 1, MOD) % MOD
                m -= v
        res = res * E.C(m, sick[0]) % MOD
        return res % MOD
