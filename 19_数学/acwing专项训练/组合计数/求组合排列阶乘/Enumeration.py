# 求组合数
# mod为素数且0<=n,k<min(mod,1e7)


from collections import Counter


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

    def lucas(self, n: int, k: int) -> int:
        if k == 0:
            return 1
        mod = self._mod
        return self.C(n % mod, k % mod) * self.lucas(n // mod, k // mod) % mod

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


if __name__ == "__main__":
    MOD = int(1e9 + 7)
    C = Enumeration(size=int(1e5), mod=MOD)

    # 统计一个字符串的 k 子序列美丽值最大的数目
    # https://leetcode.cn/contest/biweekly-contest-112/problems/count-k-subsequences-of-a-string-with-maximum-beauty/
    class Solution:
        def countKSubsequencesWithMaxBeauty(self, s: str, k: int) -> int:
            MOD = int(1e9 + 7)
            counter = Counter(s)
            if len(counter) < k:
                return 0

            freqCounter = Counter(counter.values())
            items = sorted(freqCounter.items(), key=lambda x: -x[0])
            # print(items)  # "abcdd" -> [(2, 1), (1, 3)]
            res = 1
            remain = k
            for freq, sameCount in items:
                if remain <= 0:
                    break
                min_ = min(remain, sameCount)
                res *= C.C(sameCount, min_) * pow(freq, min_, MOD)
                res %= MOD
                remain -= min_

            return res % MOD

    print(Solution().countKSubsequencesWithMaxBeauty("bcca", 2))

    # https://yukicoder.me/problems/no/117
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    T = int(input())
    C = Enumeration(10**6 + 10, 10**9 + 7)
    for _ in range(T):
        s = input()
        op = s[0]
        inner = s[2:-1]
        n, k = map(int, inner.split(","))
        if op == "C":
            print(C.C(n, k))
        elif op == "P":
            print(C.P(n, k))
        elif op == "H":
            print(C.H(n, k))
