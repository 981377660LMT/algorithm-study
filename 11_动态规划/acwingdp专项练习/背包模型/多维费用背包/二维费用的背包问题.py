# 有 N 件物品和一个容量是 V 的背包，背包能承受的最大重量是 M。
# 每件物品只能用一次。体积是 vi，重量是 mi，价值是 wi。
# 求解将哪些物品装入背包，可使物品总体积不超过背包容量，总重量不超过背包可承受的最大重量，且价值总和最大。
# 输出最大价值。

n, cap1, cap2 = map(int, input().split())

goods = []
for _ in range(n):
    cost1, cost2, score = map(int, input().split())
    goods.append([cost1, cost2, score])

dp = [[0] * (cap2 + 1) for _ in range(cap1 + 1)]

for cost1, cost2, score in goods:
    for i in range(cap1, cost1 - 1, -1):
        for j in range(cap2, cost2 - 1, -1):
            dp[i][j] = max(dp[i][j], dp[i - cost1][j - cost2] + score)


print(dp[-1][-1])

