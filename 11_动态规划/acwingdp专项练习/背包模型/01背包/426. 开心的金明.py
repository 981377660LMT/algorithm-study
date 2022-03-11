# 重要度，分为 5 等：用整数 1∼5 表示，第 5 等最重要。
# 他希望在不超过 N 元（可以等于 N 元）的前提下，使每件物品的价格与重要度的乘积的总和最大。

cap, n = map(int, input().split())
goods = []
for _ in range(n):
    cost, level = map(int, input().split())
    goods.append((cost, cost * level))

dp = [0] * (cap + 1)
for cost, score in goods:
    for i in range(cap, cost - 1, -1):
        dp[i] = max(dp[i], dp[i - cost] + score)
print(dp[-1])
