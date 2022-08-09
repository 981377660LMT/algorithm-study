# 2<=k<=n<=1e6
# 1-n号车站
# 在1停站 在n停站 中间可以选择停不停
# !不能连续停>=k个车站 求方案数 => k斐波那契数 01序列中1不能连续出现k次

# !dp[index][count] 表示前index车站连续停了count个车站的方案数
# 这里不能存index*count个状态
# 观察转移表达式
# 停车:ndp[j]=dp[j-1] (j>=1)
# 不停车:ndp[0]=dp[0]+dp[1]+...+dp[k-1]
# !无法前缀和优化dp 但是表达式转移类似队列(队列整体右移、队首加一个和)
# !所以deque来表示dp数组 记录deque里所有数的和
# !答案为dp[1]+dp[2]+...+dp[k-1] = sum_ - dp[0] (最后一站要停)

from collections import deque
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

n, k = map(int, input().split())


dp = deque([0] * k)
dp[1] = 1
sum_ = 1
for _ in range(n - 1):
    last = dp.pop()
    dp.appendleft(sum_)
    sum_ += sum_ - last
    sum_ %= MOD

print((sum_ - dp[0]) % MOD)
