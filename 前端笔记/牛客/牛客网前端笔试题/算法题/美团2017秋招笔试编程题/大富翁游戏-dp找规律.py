# 玩家根据骰子的点数决定走的步数，
# 即骰子点数为1时可以走一步，点数为2时可以走两步，
# 点数为n时可以走n步。
# 求玩家走到第n步（n<=骰子最大点数且是方法的唯一入参）时，
# 总共有多少种投骰子的方法。

# dp[n]=dp[n-1]+dp[n-2]+...+dp[1]+dp[0]
# 即dp[n]=2*dp[n-1]
n = int(input())
print(pow(2, n - 1))

