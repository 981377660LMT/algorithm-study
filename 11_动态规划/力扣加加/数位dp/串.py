# 长度不超过n，且包含子序列“us”的、只由小写字母构成的字符串有多少个？ 答案对10^9+7 取模
# n<=10^6
# 每个位置三种状态，之前 不包含u/只包含u/包含us
n = int(input())
MOD = int(1e9 + 7)
res = 0

# dp = [[0, 0, 0] for _ in range(n + 1)]

# dp[1][0] = 25
# dp[1][1] = 1
# dp[1][2] = 0

# for i in range(2, n + 1):
#     dp[i][0] = dp[i - 1][0] * 25 % MOD
#     dp[i][1] = (dp[i - 1][1] * 25 + dp[i - 1][0]) % MOD
#     dp[i][2] = (dp[i - 1][1] + dp[i - 1][2] * 26) % MOD
#     res += dp[i][2]
#     res %= MOD
# print(res)


# 滚动数组优化:
dp = [1, 0, 0]
for i in range(1, n + 1):
    dp[0], dp[1], dp[2] = (
        dp[0] * 25 % MOD,
        (dp[1] * 25 + dp[0]) % MOD,
        (dp[1] + dp[2] * 26) % MOD,
    )
    res += dp[2] % MOD
    res %= MOD
print(res)
