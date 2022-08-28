# 完美世界

n, mp = map(int, input().split())

skill = []

for _ in range(n):
    damage, cost = map(int, input().split())
    skill.append((damage, cost))


dp = [0] * (mp + 1)
for damage, cost in skill:
    for cap in range(cost, mp + 1):
        dp[cap] = max(dp[cap], dp[cap - cost] + damage)

print(dp[mp])
