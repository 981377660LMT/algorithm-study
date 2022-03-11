# 有一个箱子容量为 V，同时有 n 个物品，每个物品有一个体积（正整数）。
# 要求 n 个物品中，任取若干个装入箱内，使箱子的剩余空间为最小。

cap = int(input())
n = int(input())
goods = []
for i in range(n):
    cost = int(input())
    goods.append(cost)

dp = [False] * (cap + 1)
dp[0] = True
for cost in goods:
    for i in range(cap, cost - 1, -1):
        dp[i] = dp[i] or dp[i - cost]
print(cap - next(i for i in range(cap, -1, -1) if dp[i]))
