# W 公司有 m 个仓库和 n 个零售商店。
# 第 i 个仓库有 ai 个单位的货物；第 j 个零售商店需要 bj 个单位的货物。
# 货物供需平衡，即∑i=1mai=∑j=1nbj。
# !从第 i 个仓库运送每单位货物到第 j 个零售商店的费用为 cij。
# 试设计一个将仓库中所有货物运送到零售商店的运输方案。
# 对于给定的 m 个仓库和 n 个零售商店间运送货物的费用，计算最优运输方案和最差运输方案。
# https://www.acwing.com/problem/content/description/2194/


from MinCostMaxFlow import MinCostMaxFlowDinic

from collections import defaultdict

# 仓库数和零售商店数
# 1≤m≤100,
# 1≤n≤50,
m, n = map(int, input().split())

# 表示第 i 个仓库有 ai 个单位的货物。
stores = list(map(int, input().split()))

# 第 j 个零售商店需要 bj 个单位的货物。
needs = list(map(int, input().split()))

# 表示从第 i 个仓库运送每单位货物到第 j 个零售商店的费用
dist = defaultdict(lambda: defaultdict(lambda: int(1e20)))
for i in range(m):
    nums = list(map(int, input().split()))
    for j, num in enumerate(nums):
        dist[i][j] = num


START, END, OFFSET = 2 * (n + m + 1), 2 * (n + m + 2), n + m

# 最小费用
mcmf1 = MinCostMaxFlowDinic(2 * (n + m + 3), START, END)
for i in range(m):
    mcmf1.addEdge(START, i, stores[i], 0)  # 虚拟源点提货物
for i in range(n):
    mcmf1.addEdge(i + OFFSET, END, needs[i], 0)  # 虚拟汇点接受货物
for i in dist:
    for j in dist[i]:
        mcmf1.addEdge(i, j + OFFSET, stores[i], dist[i][j])  # 仓库转移虚拟源点的货物
print(mcmf1.work()[1])

# 最大费用
mcmf2 = MinCostMaxFlowDinic(2 * (n + m + 3), START, END)
for i in range(m):
    mcmf2.addEdge(START, i, stores[i], 0)
for i in range(n):
    mcmf2.addEdge(i + OFFSET, END, needs[i], 0)
for i in dist:
    for j in dist[i]:
        mcmf2.addEdge(i, j + OFFSET, stores[i], -dist[i][j])
print(-mcmf2.work()[1])
