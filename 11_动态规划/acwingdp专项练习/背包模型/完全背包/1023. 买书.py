# 小明手里有n元钱全部用来买书，书的价格为10元，20元，50元，100元。
# 问小明有多少种买书方案？（每种书可购买多本）
n = int(input())
dp = [0] * (n + 1)
dp[0] = 1

for num in [10, 20, 50, 100]:
    for i in range(num, n + 1):
        dp[i] += dp[i - num]

print(dp[-1])
