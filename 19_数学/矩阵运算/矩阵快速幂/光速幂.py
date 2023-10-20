# 光速幂
# 已知base和模数mod，求base^n 模mod
# !O(sqrt(maxN))预处理,O(1)查询


from math import sqrt


MOD = int(1e9 + 7)


class FastPow:
    __slots__ = "_max", "_divData", "_modData"

    def __init__(self, base: int, maxN: int) -> None:
        max_ = int(sqrt(maxN)) + 1
        self._max = max_
        self._divData = [0] * (max_ + 1)
        self._modData = [0] * (max_ + 1)
        cur = 1
        for i in range(max_ + 1):
            self._modData[i] = cur
            cur = cur * base % MOD
        cur = 1
        last = self._modData[max_]
        for i in range(max_ + 1):
            self._divData[i] = cur
            cur = cur * last % MOD

    def pow(self, n: int) -> int:
        """n<=maxN."""
        return self._divData[n // self._max] * self._modData[n % self._max] % MOD

    def rangePow2Sum(self, start: int, end: int) -> int:
        """区间以2为底的幂和 (2^start + 2^(start+1) + ... + 2^(end-1)) % MOD."""
        if start >= end:
            return 0
        return (self.pow(end) - self.pow(start)) % MOD


if __name__ == "__main__":
    fp = FastPow(2, 100)
    print(fp.pow(10))
