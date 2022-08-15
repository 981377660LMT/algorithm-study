# 高桥从1出发，想要走到n，
# 每个位置有一个面的骰子，可能等概率的摇到任意0-nums[i],询问走到点n的期望步数

# 概率dp
# !类似于新21点 后缀和倒着推

# dp[i]表示从i走到n的期望步数
# !dp[i] = dp[i] + dp[i + 1] + dp[i + 2] + ... + dp[i + nums[i]] / (nums[i] + 1)
# dp[i]全移到左边
# !dp[i] =  dp[i + 1] + dp[i + 2] + ... + dp[i + nums[i]] / nums[i] + (nums[i]+1)/(nums[i])
# 后缀优化dp O(n)


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
nums = list(map(int, input().split()))


dp, dpSum = [0] * (n + 10), [0] * (n + 10)
for i in range(n - 1, 0, -1):
    # !注意除法要变成乘逆元
    num = nums[i - 1]
    inv = pow(num, MOD - 2, MOD)
    dp[i] = (dpSum[i + 1] - dpSum[i + num + 1] + 1 + num) * inv
    dp[i] %= MOD
    dpSum[i] = dpSum[i + 1] + dp[i]
    dpSum[i] %= MOD
print(dp[1])
