"""多重背包单调队列优化 时间复杂度O(n*v*logs)"""

# 有 N 种物品和一个容量是 V 的背包。
# 第 i 种物品`最多有 si 件`，每件体积是 vi，价值是 wi。
# 求解将哪些物品装入背包，可使物品体积总和不超过背包容量，且价值总和最大。
# 输出最大价值。

# 0<N≤1000
# 0<V≤2000
# 0<vi,wi,si≤2000


# 二进制做法，对于每种物品i而言，将其拆分为Log(si) + 1种物品
# 转换后的物品每种有且只有一个，即将问题转换为 0-1 背包
# 每种新的物品的体积和价值是拆分物品的 1, 2, 4, 8... 倍

# !时间复杂度O(nlog(count))
n, cap = map(int, input().split())
goods = []
for _ in range(n):
    cost, score, count = map(int, input().split())
    cur = 1
    # 尽可能多的分成二进制组 1,2,4...
    while cur <= count:
        goods.append((cost * cur, score * cur))
        count -= cur
        cur *= 2
    # 剩下的一组
    if count:
        goods.append((cost * count, score * count))


# 转换成0-1背包
dp = [0] * (cap + 1)
for i in range(len(goods)):
    for j in range(cap, -1, -1):
        if j >= goods[i][0]:
            dp[j] = max(dp[j], dp[j - goods[i][0]] + goods[i][1])

print(dp[-1])
