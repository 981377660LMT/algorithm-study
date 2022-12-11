# n<=1e10000
# d<=100
# !不超过n的正整数中 `各数位和为d的倍数` 的数有多少个

# atc数位dp模板题


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

n = input()
d = int(input())

# !TLE
# def cal(upper: int) -> int:
#     def dfs(pos: int, isLimit: bool, sumMod: int) -> int:
#         if pos == n:
#             return sumMod == 0
#         hash = pos * (d + 1) * 2 + sumMod * 2 + isLimit
#         if memo[hash] != -1:
#             return memo[hash]
#         res = 0
#         up = nums[pos] if isLimit else 9
#         for cur in range(up + 1):
#             res += dfs(pos + 1, isLimit and cur == up, (sumMod + cur) % d)
#             res %= MOD
#         memo[hash] = res
#         return res

#     nums = list(map(int, str(upper)))
#     n = len(nums)
#     memo = [-1] * (d + 1) * (n + 1) * 2
#     return dfs(0, True, 0)


# print(cal(n) - 1)

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
