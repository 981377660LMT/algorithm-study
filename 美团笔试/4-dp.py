# 1. 球上所写的数之和可以被k1整除
# 2. 球上所写的数之和不能被k2整除
# 3. 在满足前两个条件的前提下，球上所写的数之和要尽可能的大
# 小团还想知道在满足这些条件的情况下有多少种不同的选择方法。
# n, k1 , k2 (1<=n<=100000,1<=k1 , k2<=10)

from collections import defaultdict


MOD = 998244353


n, k1, k2 = map(int, input().split())
nums = list(map(int, input().split()))
# n, k1, k2 = [5, 3, 4]
# nums = [6, 8, -2, -5, 2]
# 取模时的最大值

dp = [defaultdict(int) for _ in range(n)]
dp[0][nums[0]] = 1

# 和，模
for i in range(1, n):
    for preSum, preCount in dp[i - 1].items():
        dp[i][preSum] += preCount
        dp[i][preSum] %= MOD
        dp[i][preSum + nums[i]] += preCount
        dp[i][preSum + nums[i]] %= MOD

for maxSum in sorted(dp[-1], reverse=True):
    if (maxSum % k1 == 0) and (maxSum % k2 != 0):
        print(maxSum, dp[-1][maxSum])
        break
