# n<=1e10000
# d<=100
# !不超过n的正整数中 `各数位和为d的倍数` 的数有多少个

# 数位dp模板题 可以不用hasLeadingZero 参数

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

d = int(input())
n = int(input())


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, isLimit: bool, sumMod: int) -> int:
        if pos == n:
            return sumMod == 0

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            res += dfs(pos + 1, isLimit and cur == up, (sumMod + cur) % d)
            res %= MOD
        return res

    nums = list(map(int, str(upper)))
    n = len(nums)
    return dfs(0, True, 0)


print(cal(n) - 1)
