# # 有 N 件物品和一个容量是 V 的背包。每件物品只能使用一次。
# # 第 i 件物品的体积是 vi，价值是 wi。
# # 求解将哪些物品装入背包，可使这些物品的总体积不超过背包容量，且总价值最大。
# # 输出 字典序最小的方案。这里的字典序是指：所选物品的编号所构成的序列。物品的编号范围是 1…N。

n, cap = map(int, input().split())
goods = []
for _ in range(n):
    cost, score = map(int, input().split())
    goods.append((cost, score))

dp = [0] * (cap + 1)  # 最大总价值
dp[0] = 0
# 此处也可状压记录
assign = [[] * n for _ in range(cap + 1)]  # 每个容量下每个物品的选择个数

for i, (cost, score) in enumerate(goods):
    for j in range(cap, cost - 1, -1):
        if dp[j - cost] + score > dp[j]:
            dp[j] = dp[j - cost] + score
            assign[j] = assign[j - cost][:]
            assign[j] += [i]
        elif dp[j - cost] + score == dp[j]:
            # 字典序最小
            cand = assign[j - cost][:]
            cand += [i]
            if cand < assign[j]:
                assign[j] = cand

max_ = max(dp)
res = [int(1e20)] * n  # 字典序超大

for value, curAssign in zip(dp, assign):
    if value == max_ and curAssign < res:
        res = curAssign

for i in res:
    print(i + 1, end=' ')

