# 求组合数


class Combination:
    __slots__ = ("_mod", "_fac", "_ifac")

    def __init__(self, size: int, mod=int(1e9 + 7)) -> None:
        self._mod = mod
        self._buildFac(size)

    def fac(self, n: int) -> int:
        """阶乘"""
        return self._fac[n]

    def ifac(self, n: int) -> int:
        """阶乘逆元"""
        return self._ifac[n]

    def C(self, n: int, k: int) -> int:
        """组合数"""
        if n < 0 or k < 0 or n < k:
            return 0
        return self._fac[n] * self._ifac[k] % self._mod * self._ifac[n - k] % self._mod

    def P(self, n: int, k: int) -> int:
        """排列数"""
        if n < 0 or n < k:
            return 0
        return self._fac[n] * self._ifac[n - k] % self._mod

    def H(self, n: int, k: int) -> int:
        """
        可重复选取元素的组合数.
        `itertools.combinations_with_replacement`的元素个数.
        """
        if n == 0:
            return 1 if k == 0 else 0
        return self.C(n + k - 1, k)

    def catalan(self, n: int) -> int:
        """卡特兰数,注意2*n需要开两倍空间"""
        return self.C(2 * n, n) * pow(n + 1, self._mod - 2, self._mod) % self._mod

    def put(self, n: int, k: int) -> int:
        """n个物品放入k个槽(槽可空)的方案数"""
        return self.C(n + k - 1, k - 1)

    def _buildFac(self, n: int) -> None:
        self._fac = [1] * (n + 1)
        self._ifac = [1] * (n + 1)
        mod = self._mod
        for i in range(1, n + 1):
            self._fac[i] = self._fac[i - 1] * i % self._mod
        self._ifac[n] = pow(self._fac[n], mod - 2, mod)
        for i in range(n, 0, -1):
            self._ifac[i - 1] = self._ifac[i] * i % self._mod

    def __call__(self, n: int, k: int) -> int:
        return self.C(n, k)


if __name__ == "__main__":
    # https://yukicoder.me/problems/no/117
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    T = int(input())
    C = Combination(2 * 10**6 + 10)
    for _ in range(T):
        s = input()
        op = s[0]
        inner = s[2:-1]
        n, k = map(int, inner.split(","))
        if op == "C":
            print(C(n, k))
        elif op == "P":
            print(C.P(n, k))
        elif op == "H":
            print(C.H(n, k))
