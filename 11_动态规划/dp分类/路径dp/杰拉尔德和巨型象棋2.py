# https://atcoder.jp/contests/abc042/tasks/arc058_b
# 左上角移动到右下角，一共有多少种路线
# 每一步可以向右或向下移动一格，不可以经过障碍物
# 障碍物的范围为距离下侧范围<=A和距离左侧范围<=B的区域
# ROW<=1e6,COL<=1e6.


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


ROW, COL, A, B = map(int, input().split())


MOD = int(1e9 + 7)
E = Enumeration(int(2e5), MOD)
res = 0
for i in range(ROW - A):
    # (0,0) -> (i,B-1) -> (ROW-1,COL-1)
    a, b = i, B - 1
    cur = E.C(a + b, a)
    c, d = ROW - 1 - i, COL - B - 1
    cur = cur * E.C(c + d, c) % MOD
    res = (res + cur) % MOD
print(res % MOD)
