# 潜水员为了潜水要使用特殊的装备。
# 他有一个带2种气体的气缸：一个为氧气，一个为氮气。
# 让潜水员下潜的深度需要各种数量的氧和氮。
# 潜水员有一定数量的气缸。
# 每个气缸都有重量和气体容量。
# 潜水员为了完成他的工作需要特定数量的氧和氮。
# 他完成工作所需气缸的总重的最低限度的是多少？


# 二维背包求最小费用
cap1, cap2 = map(int, input().split())
n = int(input())
goods = []
for _ in range(n):
    cost1, cost2, score = map(int, input().split())
    goods.append((cost1, cost2, score))

dp = [[int(1e20)] * (cap2 + 1) for _ in range(cap1 + 1)]
dp[0][0] = 0

for cost1, cost2, score in goods:
    for i in range(cap1, -1, -1):
        for j in range(cap2, -1, -1):
            # 注意这里的写法 cap 为负数表示超出了需要 需要置为0
            dp[i][j] = min(dp[i][j], dp[max(0, i - cost1)][max(0, j - cost2)] + score)
print(dp[cap1][cap2])

