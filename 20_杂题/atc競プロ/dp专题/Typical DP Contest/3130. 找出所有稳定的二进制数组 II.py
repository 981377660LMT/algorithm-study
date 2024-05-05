# https://leetcode.cn/problems/find-all-possible-stable-binary-arrays-ii/description/
# 3130. 找出所有稳定的二进制数组 II
# 求二进制数组个数，满足：
# 0 在 arr 中出现次数 恰好 为 zero 。
# 1 在 arr 中出现次数 恰好 为 one 。
# 0 或1 最多连续出现limit次.
# 组合数学解法，O(zero*one/limit)
# https://leetcode.cn/problems/find-all-possible-stable-binary-arrays-ii/solutions/2758868/dong-tai-gui-hua-cong-ji-yi-hua-sou-suo-37jdi/


class Comb:
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

    def __call__(self, n: int, k: int) -> int:
        return self.C(n, k)


MOD = int(1e9 + 7)
C = Comb(int(1e5 + 10), MOD)


class Solution:
    def numberOfStableArrays(self, zero: int, one: int, limit: int) -> int:
        if zero > one:
            zero, one = one, zero  # 保证空间复杂度为 O(min(zero, one))
        f0 = [0] * (zero + 3)
        for i in range((zero - 1) // limit + 1, zero + 1):
            f0[i] = C(zero - 1, i - 1)
            for j in range(1, (zero - i) // limit + 1):
                f0[i] = (
                    f0[i] + (-1 if j % 2 else 1) * C(i, j) * C(zero - j * limit - 1, i - 1)
                ) % MOD

        res = 0
        for i in range((one - 1) // limit + 1, min(one, zero + 1) + 1):
            f1 = C(one - 1, i - 1)
            for j in range(1, (one - i) // limit + 1):
                f1 = (f1 + (-1 if j % 2 else 1) * C(i, j) * C(one - j * limit - 1, i - 1)) % MOD
            res = (res + (f0[i - 1] + f0[i] * 2 + f0[i + 1]) * f1) % MOD
        return res
