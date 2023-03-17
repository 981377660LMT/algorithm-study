from prime import EratosthenesSieve
from typing import Callable


class PrimeSum:
    """https://maspypy.github.io/library/nt/primesum.hpp

    给定n和函数f, 计算前缀和 sum_{p <= x} f(p),
    其中x必须要形如floor(n/i)的形式。

    O(n^(3/4)/logn) time, O(n^(1/2)) space.
    """

    __slots__ = ("_n", "sqN", "_sumLo", "_sumHi")

    def __init__(self, n: int) -> None:
        self._n = n
        self.sqN = int(n**0.5)

    def cal(self, f: Callable[[int], int]) -> None:
        primes = EratosthenesSieve(self.sqN).getPrimes()
        self._sumLo = [0] * (self.sqN + 1)
        self._sumHi = [0] * (self.sqN + 1)
        for i in range(1, self.sqN + 1):
            self._sumLo[i] = f(i) - 1
            self._sumHi[i] = f(self._n // i) - 1
        for p in primes:
            pp = p * p
            if pp > self._n:
                break
            R = min(self.sqN, self._n // pp)
            M = self.sqN // p
            x = self._sumLo[p - 1]
            fp = self._sumLo[p] - self._sumLo[p - 1]
            for i in range(1, M + 1):
                self._sumHi[i] -= fp * (self._sumHi[i * p] - x)
            for i in range(M + 1, R + 1):
                self._sumHi[i] -= fp * (self._sumLo[self._n // (i * p)] - x)
            for i in range(self.sqN, pp - 1, -1):
                self._sumLo[i] -= fp * (self._sumLo[i // p] - x)

    def calCount(self) -> None:
        return self.cal(lambda x: x)

    def calSum(self) -> None:
        return self.cal(lambda x: (x * (x + 1)) >> 1)

    def __getitem__(self, x: int) -> int:
        """x 必须要是 self._n//i 的形式"""
        return self._sumLo[x] if x <= self.sqN else self._sumHi[self._n // x]


if __name__ == "__main__":
    N = int(2e8)
    ps = PrimeSum(N)
    ps.calSum()
    print(ps[N])
