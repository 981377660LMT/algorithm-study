# 1<=ROW,COL<=1000
# 在 n×m 的棋盘上放置 k个棋子
# 记矩形 A 为能覆盖所有 k 个棋子的最小的矩形
# 求 A 的面积的期望

# !0.算期望等于所有可能的值除以总的可能方案数.
# !1.枚举每种矩形有多少种防止棋子的方案.
# !2.对于一个 i×j 的矩形，我们可以用容斥的方法求出正好覆盖它的方案数.


MOD = 998244353


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


E = Enumeration(int(1e6), MOD)


def minimumBoundingBox2(ROW: int, COL: int, k: int) -> int:
    def f(i: int, j: int) -> int:
        """
        边长为i*j的矩形,有多少种选出k个格子的方案,且这个矩形为最小矩形.
        返回选出的所有格子的个数。

        减去四种情况:
        - 第一行没有棋子
        - 最后一行没有棋子
        - 第一列没有棋子
        - 最后一列没有棋子
        """
        res = E.C(i * j, k)
        if i > 1:
            res -= 2 * E.C((i - 1) * j, k)
        if j > 1:
            res -= 2 * E.C(i * (j - 1), k)
        if i > 1 and j > 1:
            res += 4 * E.C((i - 1) * (j - 1), k)
        if i > 2:
            res += E.C((i - 2) * j, k)
        if j > 2:
            res += E.C(i * (j - 2), k)
        if i > 2 and j > 1:
            res -= 2 * E.C((i - 2) * (j - 1), k)
        if i > 1 and j > 2:
            res -= 2 * E.C((i - 1) * (j - 2), k)
        if i > 2 and j > 2:
            res += E.C((i - 2) * (j - 2), k)
        res %= MOD
        cand = i * j  # !选出的格子数
        count = (ROW - i + 1) * (COL - j + 1)  # !有多少个这样的矩形
        return res * cand % MOD * count % MOD

    res = 0
    for i in range(1, ROW + 1):
        for j in range(1, COL + 1):
            if i * j < k:
                continue
            res = (res + f(i, j)) % MOD

    all_ = E.C(ROW * COL, k)
    inv = pow(all_, MOD - 2, MOD)
    return res * inv % MOD


if __name__ == "__main__":
    ROW, COL, k = map(int, input().split())
    print(minimumBoundingBox2(ROW, COL, k))
