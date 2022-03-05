# 一个正整数 n 可以表示成若干个正整数之和，形如：n=n1+n2+…+nk，
# 其中 n1≥n2≥…≥nk,k≥1。
# 我们将这样的一种表示称为正整数 n 的一种划分。
# 现在给定一个正整数 n，请你求出 n 共有多少种不同的划分方法。

# 1≤n≤1000

# n个物品，第i个物品的体积和价值都为i，并且求价值恰好为n的划分方案个数
MOD = int(1e9 + 7)
n = int(input().strip())
dp = [0] * (n + 1)
dp[0] = 1  # 什么都不选

# 组物外
for i in range(1, n + 1):
    for j in range(1, n + 1):
        if j - i >= 0:
            dp[j] += dp[j - i]
            dp[j] %= MOD

print(dp[-1])
