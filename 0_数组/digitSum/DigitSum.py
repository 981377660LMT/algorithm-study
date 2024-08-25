# !预处理更快

import time


class DigitSum:
    __slots__ = ("_mod", "_dp")

    def __init__(self, step: int = 6) -> None:
        step = max(4, min(7, step))
        self._mod = 10**step
        self._dp = [0] * self._mod
        for x in range(1, self._mod):
            self._dp[x] = self._dp[x // 10] + (x % 10)

    def sum(self, x: int) -> int:
        res = 0
        dp, mod = self._dp, self._mod
        while x > 0:
            res += dp[x % mod]
            x //= mod
        return res


def digit_sum_naive(num: int) -> int:
    sum = 0
    while num > 0:
        sum += num % 10
        num //= 10
    return sum


if __name__ == "__main__":
    time1 = time.time()
    ds = DigitSum()
    for i in range(10**7):
        ds.sum(i)
    time2 = time.time()
    print(time2 - time1)  # 4.37386417388916

    time1 = time.time()
    for i in range(10**7):
        digit_sum_naive(i)
    time2 = time.time()
    print(time2 - time1)  # 8.102303504943848
