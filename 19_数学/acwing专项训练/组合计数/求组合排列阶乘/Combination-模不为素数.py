"""https://atcoder.jp/users/zatsua"""
# 带模的组合数计算


from typing import Tuple


class Factorial:
    @staticmethod
    def xgcd(a: int, b: int) -> Tuple[int, int, int]:
        """return (g, x, y) such that a*x + b*y = g = gcd(a, b) ;扩展欧几里得"""
        x0, x1, y0, y1 = 0, 1, 1, 0
        while a != 0:
            (q, a), b = divmod(b, a), a
            y0, y1 = y1, y0 - q * y1
            x0, x1 = x1, x0 - q * x1
        return b, x0, y0

    def __init__(self, MOD: int):
        self._mod = MOD
        self._fac = [1]
        self._size = 1
        self._iFac = [1]
        self._iSize = 1

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

    def P(self, n: int, k: int) -> int:
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

    def Catalan(self, n: int) -> int:
        """卡特兰数"""
        return self.C(2 * n, n) * self._modinv(n + 1) % self._mod

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


if __name__ == "__main__":
    f = Factorial(int(1e9 + 7))
    print(f.CWithReplacement(5, 2))
    print(f.CWithReplacement(4, 2))
    print(f.put(3, 2))  # 4
    print(f.put(4, 2))  # 5

    for i in range(5):
        print(f.Catalan(i))
