#  (1 ≤n, m, k < 50)
# !物品总数<=998244353
# !O(n*m*k^2)

# !n种物品 每种物品有无限个且每个物品获得的概率为wi/∑(wi)
# !有放回地选k次（每次选一个物品）求有m个不同物品的可能性
# !每次选择的概率独立 (mod 998244353)

# 组合数+概率dp
# dp[i][j][k] 表示前i个彩票取了j种 一共取了k张的概率
# https://www.zhihu.com/search?type=content&q=AtCoder%20Beginner%20Contest%20243
# https://www.zhihu.com/search?q=AtCoder%20Beginner%20Contest%20243&utm_content=search_history&type=content
# TODO


import sys
import os
from typing import Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    n, m, k = map(int, input().split())
    w = [int(input()) for _ in range(n)]  # 每种物品获得的可能性

    # 遍历每个物品堆 看选取几个
    def dfs(index: int, remainCount: int, remainType: int) -> int:
        if remainCount < 0 or remainType < 0:
            return 0
        if index == n:
            return int(remainCount == 0 and remainType == 0)
        res = dfs(index + 1, remainCount, remainType)

        # 这一堆选出select个
        for select in range(1, remainCount + 1):
            # select个物品拿出来的概率 pi^select/(select!)
            prob = pow(w[index], select, MOD) * pow(F(select), MOD - 2, MOD)
            res += dfs(index + 1, remainCount - select, remainType - 1) * prob
            res %= MOD
        return res % MOD

    print(dfs(0, k, m) % MOD)


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
