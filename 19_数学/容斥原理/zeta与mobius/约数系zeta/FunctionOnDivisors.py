# TODO ？？？
# https://maspypy.github.io/library/nt/function_on_divisors.hpp


from typing import Callable
from Factor import Factor


class FunctionOnDivisors:
    __slots__ = ("_pf", "_divs", "_data")

    def __init__(self, n: int) -> None:
        pf = Factor.getPrimeFactors(n)
        pf = sorted(pf.items())
        divs = [1]
        for p, e in pf:
            n = len(divs)
            q = p
            for _ in range(e):
                for i in range(n):
                    divs.append(divs[i] * q)
                q *= p
        self._pf = pf
        self._divs = divs

    def setMultiplicative(self, f: Callable[[int, int], int]) -> None:
        """设定f(p, k)"""
        data = [1]
        n = len(self._divs)
        for p, e in self._pf:
            for k in range(1, e + 1):
                for i in range(n):
                    data.append(data[i] * f(p, k))
        self._data = data

    def setEulerPhi(self) -> None:
        """设定欧拉函数"""
        self._data = self._divs[:]
        self.divisorMobius()

    def enumerate(self, f: Callable[[int, int], None]) -> None:
        for d, fd in zip(self._divs, self._data):
            f(d, fd)

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


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc212/tasks/abc212_g
    def f(p: int, k: int) -> None:
        global res
        print(p, k)
        res += p * k

    p = int(input())
    F = FunctionOnDivisors(p - 1)
    F.setEulerPhi()
    res = 1
    F.enumerate(f)
    print(res % 998244353)
