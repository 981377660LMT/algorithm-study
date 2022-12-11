# 给定n个点 每次随机选择1个点移动
# 求使图连通的移动次数的期望值
# dp[i]表示已经有i个点连通
# !dp[i] = 1 + (i/n)*dp[i] + (1-i/n)*dp[i+1]
# !即 dp[i] = dp[i+1] + n/(n-i),dp[n]=0


def journey(n: int) -> float:
    dp = 0
    for i in range(n - 1, 0, -1):
        dp += n / (n - i)
    return dp


print(journey(int(input())))
