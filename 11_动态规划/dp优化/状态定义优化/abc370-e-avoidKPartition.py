# E - Avoid K Partition (正难则反)
# https://atcoder.jp/contests/abc370/tasks/abc370_e
# 分割成k个子数组，不含和为K的方案数，模998244353
# n<=2e5
# dp[i] 表示最后一个分割点为i的方案数


import sys
from collections import defaultdict


input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))

    mp = defaultdict(int, {0: 1})
    curSum = 0
    all_ = 1
    dp = [0] * n  # !dp[i] 表示最后一个分割点为nums[i]之前的方案数
    for i, v in enumerate(nums):
        curSum += v
        ban = curSum - k
        cur = (all_ - mp[ban]) % MOD
        dp[i] = cur
        all_ = (all_ + cur) % MOD
        mp[curSum] = (mp[curSum] + cur) % MOD
    print(dp[-1])
