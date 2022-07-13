# AtCoder Beginner Contest 259 - SGColin的文章 - 知乎
# https://zhuanlan.zhihu.com/p/539701972

# 给定一个矩阵 Anxn (1≤n ≤400)，从某个格子出发，每次可以向右或向下走。
# !问起点终点的数字相同的路径有多少条?

# 分情况
# 暴力+dp两种算法的结合

# !注意:product/combinations会比正常循环慢一些
# !不要乱用for循环

from collections import defaultdict
from itertools import combinations
import sys
import os
from typing import List, Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353

Point = Tuple[int, int]
DIR2 = [[0, 1], [1, 0]]


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


def main() -> None:
    def strategy1(points: List[Point]) -> int:
        """brute force 枚举起始点"""
        res = len(points)  # 自己走到自己
        for (sr, sc), (er, ec) in combinations(points, 2):
            if er >= sr and ec >= sc:
                res += F.C(er - sr + ec - sc, ec - sc)
                res %= MOD
        return res

    def strategy2(points: List[Point]) -> int:
        """dp (当前横坐标，当前纵坐标) 为终点的路径数"""

        dp = [[0] * n for _ in range(n)]
        for r, c in points:
            dp[r][c] = 1
        for r in range(n):
            for c in range(n):
                # !注意这里不要用 for 循环
                if r + 1 < n:
                    dp[r + 1][c] += dp[r][c]
                    dp[r + 1][c] %= MOD
                if c + 1 < n:
                    dp[r][c + 1] += dp[r][c]
                    dp[r][c + 1] %= MOD
        res = 0
        for r, c in points:
            res += dp[r][c]
            res %= MOD
        return res

    n = int(input())
    matrix = []
    counter = defaultdict(list)
    for r in range(n):
        row = tuple(map(int, input().split()))
        for c, num in enumerate(row):
            counter[num].append((r, c))
        matrix.append(row)

    res = 0
    for points in counter.values():
        if len(points) <= n:
            res += strategy1(points)  # !枚举+组合数 最多O(n^3)的计算量
            res %= MOD
        else:
            res += strategy2(points)  # !dp 最多O(n^3)的计算量 因为最多n种数字
            res %= MOD

    print(res)


if __name__ == "__main__":

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
