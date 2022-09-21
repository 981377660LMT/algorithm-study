# Predilection (偏爱)
# 给定一个长度为N的序列，你可以执行0次或者多次一下操作。
# !选择序列中的两个相邻的数，删除这两个，并用他们的和来插入到他们原来的位置(就是合并)
# 问最后可以产生多少种不同的序列。
# 1<=n<=2e5

# https://zhuanlan.zhihu.com/p/460595293
# !dp[i+1]=2*dp[i]-d[j] 其中 sum(dp[j+1:i+1])=0
# 在减去区间和为0的段时，需要减去最近的 j


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    dp = [0] * (n + 1)
    dp[1] = 1
    preSum, curSum = dict(), 0
    for i in range(1, n):
        curSum += nums[i - 1]
        preIndex = preSum.get(curSum, 0)
        dp[i + 1] = (2 * dp[i] - dp[preIndex]) % MOD
        preSum[curSum] = i
    print(dp[n])
