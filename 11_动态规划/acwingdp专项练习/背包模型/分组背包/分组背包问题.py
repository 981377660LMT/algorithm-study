# 有 N 组物品和一个容量是 V 的背包。
# !每组物品有若干个，同一组内的物品最多只能选一个。
# 每件物品的体积是 vij，价值是 wij，其中 i 是组号，j 是组内编号。
# 求解将哪些物品装入背包，可使物品总体积不超过背包容量，且总价值最大。

# 0<N,V≤100
# 0<Si≤100
# 0<vij,wij≤100

# 和01背包区别就是枚举容量后多了一维，还要枚举编号
n, cap = map(int, input().split())
groups = [[] for _ in range(n)]
groupCounts = []
for i in range(n):
    count = int(input())
    groupCounts.append(count)
    for _ in range(count):
        cost, score = map(int, input().split())
        groups[i].append((cost, score))

dp = [0] * (cap + 1)
for i1 in range(n):
    for i2 in range(cap, -1, -1):
        tmp = dp[i2]
        for i3 in range(groupCounts[i1]):
            if i2 >= groups[i1][i3][0]:
                tmp = max(tmp, dp[i2 - groups[i1][i3][0]] + groups[i1][i3][1])
        dp[i2] = tmp
print(dp[-1])

