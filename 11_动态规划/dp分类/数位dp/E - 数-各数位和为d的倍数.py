# n<=1e10000
# d<=100
# !不超过n的正整数中 `各数位和为d的倍数` 的数有多少个

# atc数位dp模板题


from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

n = input()
d = int(input())

# !TLE


# def cal(upper: int, mod: int) -> int:
#     @lru_cache(None)
#     def dfs(pos: int, isLimit: bool, curMod: int) -> int:
#         if pos == n:
#             return curMod == 0
#         res = 0
#         up = nums[pos] if isLimit else 9
#         for cur in range(up + 1):
#             res += dfs(pos + 1, isLimit and cur == up, (curMod + cur) % mod)
#             res %= MOD
#         return res

#     nums = list(map(int, str(upper)))
#     n = len(nums)
#     res = dfs(0, True, 0)
#     dfs.cache_clear()
#     return res

# print(cal(n, d) - cal(0, d))


dp = [[0, 0] for _ in range(d)]  # sumMod isLimit
dp[0][1] = 1
for pos, char in enumerate(n):
    num = int(char)
    ndp = [[0, 0] for _ in range(d)]
    for mod in range(d):
        for isLimit in range(2):
            upper = num if isLimit else 9
            for cur in range(upper + 1):
                ndp[(mod + cur) % d][isLimit & (cur == upper)] += dp[mod][isLimit]
                ndp[(mod + cur) % d][isLimit & (cur == upper)] %= MOD
    dp = ndp

print((dp[0][0] + dp[0][1] - 1) % MOD)
