# 有 N 件物品和一个容量是 V 的背包。每件物品只能使用一次。
# 第 i 件物品的体积是 vi，价值是 wi。
# 求解将哪些物品装入背包，可使这些物品的总体积不超过背包容量，且总价值最大。
# 输出 最优选法的方案数。注意答案可能很大，请输出答案模 109+7 的结果。

MOD = int(1e9 + 7)
n, cap = map(int, input().split())
goods = []
for _ in range(n):
    cost, score = map(int, input().split())
    goods.append((cost, score))


dp = [-int(1e20)] * (cap + 1)  # 最大总价值
dp[0] = 0
count = [0] * (cap + 1)  # 取到最大总价值的最优方案个数
count[0] = 1  # 表示价值为0, 方案个数为1
for cost, score in goods:
    for i in range(cap, cost - 1, -1):
        if dp[i - cost] + score > dp[i]:
            dp[i] = dp[i - cost] + score
            count[i] = count[i - cost]
        elif dp[i - cost] + score == dp[i]:
            count[i] += count[i - cost]

max_ = max(dp)
res = 0
for v, c in zip(dp, count):
    if v == max_:
        res += c
        res %= MOD

print(res)

