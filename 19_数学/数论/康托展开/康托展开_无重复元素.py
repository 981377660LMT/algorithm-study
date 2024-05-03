"""
康托展开 - 无重复元素
求字典序第k小的排列/当前排列在所有排列中的字典序第几小

注意:
如果问比当前排列大/小 k的排列, 并不能先展开然后再放回.
!取模会丢失信息, 应该直接计算.
"""
# https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
# ! 取模的情况, len(s)很大


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
E = Enumeration(10**5 + 10, MOD)


def calRank(perm: List[int]) -> int:
    """求当前排列在所有排列中的字典序第几小(rank>=0)"""

    def add(i: int, val: int):
        while i <= n:
            bit[i] += val
            i += i & -i

    def preSum(i: int) -> int:
        res = 0
        while i > 0:
            res += bit[i]
            i &= i - 1
        return res

    n = len(perm)
    bit = [0] * (n + 1)
    for i in range(1, n + 1):
        add(i, 1)
    res = 0
    for i, v in enumerate(perm):
        res += preSum(v - 1) * E.fac(n - 1 - i) % MOD
        res %= MOD
        add(v, -1)
    return res  # 从0开始的排名


def calPerm(n: int, rank: int) -> List[int]:
    """求在1-n的所有排列中,字典序第几小(rank>=0)是谁.kthPermutation"""
    fac = [1] * (n + 10)
    for i in range(1, n + 10):
        fac[i] = fac[i - 1] * i
    perm = [0] * n
    valid = [True] * (n + 1)
    for i in range(1, n + 1):
        order = rank // fac[n - i] + 1
        for j in range(1, n + 1):
            order -= valid[j]
            if order == 0:
                perm[i - 1] = j
                valid[j] = False
                break
        rank %= fac[n - i]
    return perm


if __name__ == "__main__":
    print(calRank([1, 2, 3, 4, 5, 6, 7, 8, 10, 9]))
    print(calPerm(10, 1))

    class Solution:
        def getPermutationIndex(self, perm: List[int]) -> int:
            return calRank(perm)
