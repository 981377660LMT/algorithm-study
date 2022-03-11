# 给你一个n种面值的货币系统，求组成面值为m的货币有多少种方案。
n, m = map(int, input().split())
goods = []
for i in range(n):
    goods.append(int(input()))
dp = [0] * (m + 1)
dp[0] = 1

for num in goods:
    for i in range(num, m + 1):
        dp[i] = dp[i] + dp[i - num]
print(dp[-1])
