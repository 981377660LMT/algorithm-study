# 给定 N 个正整数 A1,A2,…,AN，从中选出若干个数，使它们的和为 M，求有多少种选择方案。
n, m = map(int, input().split())
nums = list(map(int, input().split()))

dp = [0] * (m + 1)
dp[0] = 1
for num in nums:
    for i in range(m, num - 1, -1):
        dp[i] = dp[i] + dp[i - num]
print(dp[m])
