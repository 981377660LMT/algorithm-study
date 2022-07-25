# AtCoder Beginner Contest 242 F(组合数学) - 严格鸽的文章 - 知乎
# https://zhuanlan.zhihu.com/p/476659976

# 有一个n * m(n, m ≤50)的象棋盘。
# 你需要在棋盘上放置B个黑棋， W个白棋，
# 问有多少种放置方案满足两个不同颜色的棋不在同一行或者同一列。

# !1.枚举白棋占了i行j列
# 则黑棋可以放置的格子个数为 n*m - i*m -j*n + i*j 个 (i*j重复减去了)
# 当前的 [i,j] 贡献的答案为
# C(n, i) * C(m, j) * C(n * m - i * m - j * n + i * j, BLACK)*(白棋占了i行j列)

# !2.白旗占了i行j列方案数怎么求 (dp+组合数求)
# 很明显，有些放置方案是占满不了 i 行 j 列
# 需要枚举i*j中有 x行没被用 y列没被用
# !则白旗可以放置的格子个数为 C(i*j, WHITE) - C(i,x) * C(j,y)*(白旗占了(i-x)行(j-y)列)

import sys
import os
from typing import Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:

    ROW, COL, BLACK, WHITE = map(int, input().split())

    res = 0
    dp = [[0] * (COL + 1) for _ in range(ROW + 1)]  # !白棋占了i行j列的方案数

    for i in range(1, ROW + 1):
        for j in range(1, COL + 1):
            if i * j < WHITE:
                continue

            dp[i][j] = F.C(i * j, WHITE)
            for x in range(i):
                for y in range(j):
                    if x == 0 and y == 0:
                        continue
                    dp[i][j] -= F.C(i, x) * F.C(j, y) * dp[i - x][j - y]
                    dp[i][j] %= MOD

            res += (
                F.C(ROW, i)
                * F.C(COL, j)
                * dp[i][j]
                * F.C(ROW * COL - ROW * j - COL * i + i * j, BLACK)
            )
            res %= MOD

    print(res)


if __name__ == "__main__":

    class Factorial:
        def __init__(self, MOD: int):
            self._mod = MOD
            self._fac = [1]
            self._size = 1
            self._iFac = [1]
            self._iSize = 1

        @staticmethod
        def xgcd(a: int, b: int) -> Tuple[int, int, int]:
            """return (g, x, y) such that a*x + b*y = g = gcd(a, b) ;扩展欧几里得"""
            x0, x1, y0, y1 = 0, 1, 1, 0
            while a != 0:
                (q, a), b = divmod(b, a), a
                y0, y1 = y1, y0 - q * y1
                x0, x1 = x1, x0 - q * x1
            return b, x0, y0

        def F(self, n: int) -> int:
            """n! % mod; factorial"""
            if n >= self._mod:
                return 0
            self._make(n)
            return self._fac[n]

        def C(self, n: int, k: int) -> int:
            """nCk % mod; combination"""
            if n < 0 or k < 0 or n < k:
                return 0
            t = self._fact_inv(n - k) * self._fact_inv(k) % self._mod
            return self(n) * t % self._mod

        def A(self, n: int, k: int) -> int:
            """nPk % mod"""
            if n < 0 or k < 0 or n < k:
                return 0
            return self(n) * self._fact_inv(n - k) % self._mod

        def CWithReplacement(self, n: int, k: int) -> int:
            """nHk % mod
            可重复选取的组合数 itertools.combinations_with_replacement 的个数
            """
            t = self._fact_inv(n - 1) * self._fact_inv(k) % self._mod
            return self(n + k - 1) * t % self._mod

        def put(self, n: int, k: int) -> int:
            """n个物品放入k个槽(槽可空)的方案数"""
            return self.C(n + k - 1, k - 1)

        def _fact_inv(self, n: int) -> int:
            """n!^-1 % mod"""
            if n >= self._mod:
                raise ValueError("Modinv is not exist! arg={}".format(n))
            self._make_inv(n)
            return self._iFac[n]

        # modinv(a)はax≡1(modp)となるxをreturnする。
        # ax≡y(modp)となるxは上のreturnのy倍
        def _modinv(self, n: int) -> int:
            gcd_, x, _ = self.xgcd(n, self._mod)
            if gcd_ != 1:
                raise ValueError("Modinv is not exist! arg={}".format(n))
            return x % self._mod

        def _make(self, n: int) -> None:
            mod, fac = self._mod, self._fac
            if n >= mod:
                n = mod
            if self._size < n + 1:
                for i in range(self._size, n + 1):
                    fac.append(fac[i - 1] * i % mod)
                self._size = n + 1

        def _make_inv(self, n: int) -> None:
            iFac, fac = self._iFac, self._fac
            if n >= self._mod:
                n = self._mod
            self._make(n)
            if self._iSize < n + 1:
                for i in range(self._iSize, n + 1):
                    iFac.append(self._modinv(fac[i]))
                self._iSize = n + 1

        def __call__(self, n: int) -> int:
            """n! % mod"""
            return self.F(n)

    F = Factorial(MOD)
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
