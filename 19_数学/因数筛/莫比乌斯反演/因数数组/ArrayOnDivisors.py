# https://maspypy.github.io/library/nt/function_on_divisors.hpp


from Factor import Factor
from typing import Callable, Dict


class ArrayOnDivisors:
    """维护数组的因数分解."""

    __slots__ = ("_pf", "_divs", "_data", "_mp")

    def __init__(self, pf: Dict[int, int]) -> None:
        pfList = sorted(pf.items())
        n = 1
        for _, e in pfList:
            n *= e + 1
        divs = [1] * n
        data = [0] * n
        ptr = 1
        for p, e in pfList:
            tmp = ptr
            q = p
            for _ in range(e):
                for i in range(tmp):
                    divs[ptr] = divs[i] * q
                    ptr += 1
                q *= p
        self._data = data
        self._pf = pfList
        self._divs = divs
        self._mp = {d: i for i, d in enumerate(divs)}

    def setMultiplicative(self, f: Callable[[int, int], int]) -> None:
        """设定f(p, k)"""
        data = [1]
        n = len(self._divs)
        for p, c in self._pf:
            for k in range(1, c + 1):
                for i in range(n):
                    data.append(data[i] * f(p, k))
        self._data = data

    def setEulerPhi(self) -> None:
        """设定欧拉函数"""
        self._data = self._divs[:]
        self.divisorMobius()

    def setMobius(self) -> None:
        """设定莫比乌斯函数"""
        self.setMultiplicative(lambda _, k: -1 if k == 1 else 0)

    def multiplierZeta(self) -> None:
        """倍数Zeta变换"""
        k = 1
        n = len(self._divs)
        for _, e in self._pf:
            mod = k * (e + 1)
            for i in range(n // mod):
                for j in range(mod - k - 1, -1, -1):
                    self._data[mod * i + j] += self._data[mod * i + j + k]  # op
            k *= e + 1

    def multiplierMobius(self) -> None:
        """倍数Mobius变换"""
        k = 1
        n = len(self._divs)
        for _, e in self._pf:
            mod = k * (e + 1)
            for i in range(n // mod):
                for j in range(mod - k):
                    self._data[mod * i + j] -= self._data[mod * i + j + k]  # inv
            k *= e + 1

    def divisorZeta(self) -> None:
        """约数Zeta变换"""
        k = 1
        n = len(self._divs)
        for _, e in self._pf:
            mod = k * (e + 1)
            for i in range(n // mod):
                for j in range(mod - k):
                    self._data[mod * i + j + k] += self._data[mod * i + j]
            k *= e + 1

    def divisorMobius(self) -> None:
        """约数Mobius变换"""
        k = 1
        n = len(self._divs)
        for _, e in self._pf:
            mod = k * (e + 1)
            for i in range(n // mod):
                for j in range(mod - k - 1, -1, -1):
                    self._data[mod * i + j + k] -= self._data[mod * i + j]
            k *= e + 1

    def enumerate(self, f: Callable[[int, int], None]) -> None:
        for d, fd in zip(self._divs, self._data):
            f(d, fd)

    def __getitem__(self, d: int) -> int:
        return self._data[self._mp[d]]

    def __setitem__(self, d: int, v: int) -> None:
        self._data[self._mp[d]] = v


if __name__ == "__main__":
    # # https://atcoder.jp/contests/abc212/tasks/abc212_g
    # def f(p: int, k: int) -> None:
    #     global res
    #     # print(p, k)
    #     res += p * k

    # p = int(input())
    # pf = Factor.getPrimeFactors(p - 1)
    # F = ArrayOnDivisors(pf)
    # F.setEulerPhi()
    # res = 1
    # F.enumerate(f)
    # print(res % 998244353)

    MOD = 998244353
    N, M = map(int, input().split())
    A = list(map(int, input().split()))
    pf = Factor.getPrimeFactors(M)
    divCounter = ArrayOnDivisors(pf)
    for x in A:
        if M % x == 0:
            divCounter[x] += 1
    divCounter.divisorZeta()
    dp = ArrayOnDivisors(pf)

    def f(d: int, cnt: int) -> None:
        dp[d] = pow(2, cnt, MOD) - 1

    divCounter.enumerate(f)
    dp.divisorMobius()
    print(dp[M] % MOD)
