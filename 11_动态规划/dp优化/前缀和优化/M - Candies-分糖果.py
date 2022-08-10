# 分k个糖果 每个小孩 0-ai棵
# 求方案数
# n<=100
# k<=1e5
# ai<=k

# 如果按dfs(index,remain)写 `再枚举ai个转移求和` 会TLE O(n*k^2)
# 前缀和优化dp 这个转移求和可以前缀和算出来


from itertools import accumulate
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

n, k = map(int, input().split())
nums = list(map(int, input().split()))

dp = [0] * (k + 1)  # 用i个糖果时分配方案数
for i in range(nums[0] + 1):
    dp[i] = 1

for i in range(1, n):
    dpSum = [0] + list(accumulate(dp))
    ndp = [0] * (k + 1)
    for cur in range(k + 1):
        # ndp[cur] = sum(dp[cur - ai] for ai in range(min(cur, nums[i]) + 1))
        # dp[cur] + dp[cur - 1] + ... + dp[cur - nums[i]]
        ndp[cur] = dpSum[cur + 1] - dpSum[max(0, cur - nums[i])]
        ndp[cur] %= MOD
    dp = ndp

print(dp[k])
