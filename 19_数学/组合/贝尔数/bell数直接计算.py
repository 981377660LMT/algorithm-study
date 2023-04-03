# 贝尔数B(n,k):
# n个不同的球放到不超过k个相同的盒子里的方案数
# B(n,n)表示将n个球分成任意组的方案数


class _S:
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


MOD = int(1e9 + 7)
E = _S(int(1e5), MOD)


def bell(n: int, k: int) -> int:
    """贝尔数.O(min(n,k)*logn)."""
    if k > n:
        k = n
    jsum = [0] * (k + 2)
    for j in range(k + 1):
        add = E.ifac(j)
        if j & 1:
            jsum[j + 1] = (jsum[j] - add) % MOD
        else:
            jsum[j + 1] = (jsum[j] + add) % MOD
    res = 0
    for i in range(k + 1):
        res += pow(i, n, MOD) * E.ifac(i) % MOD * jsum[k - i + 1]
        res %= MOD
    return res


if __name__ == "__main__":
    n, k = map(int, input().split())
    print(bell(n, k))
