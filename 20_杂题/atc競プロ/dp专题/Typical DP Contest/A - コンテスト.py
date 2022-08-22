import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# !A - コンテスト 01背包问题求方案数


n = int(input())
nums = list(map(int, input().split()))
sum_ = sum(nums)
dp = [False] * (sum_ + 1)
dp[0] = True
for i in range(n):
    for j in range(sum_, nums[i] - 1, -1):
        dp[j] |= dp[j - nums[i]]

print(sum(dp[i] for i in range(sum_ + 1)))

#####################################################
n = int(input())
nums = list(map(int, input().split()))
sum_ = sum(nums)
dp = [False] * (sum_ + 1)
dp[0] = True
for i in range(n):
    ndp = dp[:]
    for pre in range(sum_ + 1):
        cur = pre + nums[i]
        if cur > sum_:
            break
        ndp[cur] |= dp[pre]
    dp = ndp

print(sum(1 for i in range(sum_ + 1) if dp[i]))

#####################################################
n = int(input())
dp = 1
for num in map(int, input().split()):
    dp |= dp << num
print(bin(dp).count("1"))
