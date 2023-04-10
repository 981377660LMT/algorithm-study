# 枚举所有可能的边长i,j，先确认有多少个这样的子矩形，
# 再看矩形里靠k个填上的有多少种方法——我们使用容斥判断四个条件就好：第一行有没有，最后一行有没有，第一列有没有，最后一列有没有

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 縦
# H 行、横
# W 列のグリッドがあります。

# このグリッドから一様ランダムに
# K 個のマスを選びます。選んだマスを全て含むような（グリッドの軸に辺が平行な）最小の長方形に含まれるマスの個数がスコアとなります。


# 得られるスコアの期待値を
# mod 998244353 で求めてください。


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
        return self.C(n + k - 1, k)

    def put(self, n: int, k: int) -> int:
        """n个相同的球放入k个不同的盒子(盒子可放任意个球)的方案数."""
        return self.C(n + k - 1, n)

    def catalan(self, n: int) -> int:
        """卡特兰数"""
        return self.C(2 * n, n) * self.inv(n + 1) % self._mod

    def _expand(self, size: int) -> None:
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


E = Enumeration(int(1e6 + 10), MOD)
if __name__ == "__main__":
    ROW, COL, K = map(int, input().split())
    if K == 1:
        print(1)
        exit(0)

    all_ = E.C(ROW * COL, K)

    # 枚举被选中矩形的对角线,看有多少种命中内部选k个点的方法

    res = 1
    print(res * E.inv(all_) % MOD)
