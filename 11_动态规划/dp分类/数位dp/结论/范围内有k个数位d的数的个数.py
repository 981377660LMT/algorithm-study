# [0,upper]范围内有k个数位d的数的个数


from functools import lru_cache
from typing import List


def toString(n: int, base: int) -> List[int]:
    """Returns the base representation of n as a string.

    >>> toString(10, 2)
    [1, 0, 1, 0]
    """
    if n == 0:
        return []
    res = []
    while n:
        res.append(n % base)
        n //= base
    return res[::-1]


def calc(upper: int, k: int, d: int, *, base: int, mod: int) -> int:
    """范围[0, upper]内有k个数位d的数的个数.
    O(log(upper) * k * base) time.
    """

    assert base > 1

    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, count: int) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界,出现次数为count"""
        if count > k:
            return 0
        if pos == n:
            return count == k

        res = 0
        up = nums[pos] if isLimit else base - 1
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), count)
            else:
                res += dfs(pos + 1, False, (isLimit and cur == up), count + (cur == d))
        return res % mod

    nums = toString(upper, base)
    n = len(nums)
    res = dfs(0, True, True, 0)
    dfs.cache_clear()
    return res


if __name__ == "__main__":

    def checkWithBruteForce(upper: int, k: int, d: int, base: int) -> None:
        def bruteForce(upper: int, k: int, d: int, base: int) -> int:
            res = 0
            for i in range(upper + 1):
                if toString(i, base).count(d) == k:
                    res += 1
            return res

        assert calc(upper, k, d, base=base, mod=MOD) == bruteForce(upper, k, d, base)

    MOD = int(1e9 + 7)
    for _ in range(1000):
        from random import randint

        upper = randint(100, int(1e4))
        k = randint(0, 10)
        d = randint(0, 9)
        base = randint(2, 10)
        checkWithBruteForce(upper, k, d, base)

    # 3352. 统计小于 N 的 K 可约简整数
    # https://leetcode.cn/problems/count-k-reducible-numbers-less-than-n/
    class Solution:
        def countKReducibleNumbers(self, s: str, k: int) -> int:
            n = len(s)
            steps = [0] * (n + 1)
            for i in range(1, n + 1):
                steps[i] = steps[i.bit_count()] + 1
            res = 0
            for i in range(1, n + 1):
                if steps[i] <= k:
                    res += calc(int(s), i, 1, base=2, mod=MOD)
            return res % MOD
