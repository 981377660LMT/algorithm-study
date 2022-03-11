cap, n = map(int, input().split())
goods = []
for i in range(n):
    cost, score = map(int, input().split())
    goods.append((cost, score))

dp = [0] * (cap + 1)
for cost, score in goods:
    for i in range(cap, cost - 1, -1):
        dp[i] = max(dp[i], dp[i - cost] + score)
print(dp[-1])
