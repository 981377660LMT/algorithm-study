# G 公司有 n 个沿铁路运输线`环形`排列的仓库，每个仓库存储的货物数量不等。
# !如何用最少搬运量(相邻边的费用之和)可以使 n 个仓库的库存数量相同。
# 搬运货物时，只能在相邻的仓库之间搬运。
# 数据保证一定有解。

# 1≤n≤100,
# 每个仓库的库存量不超过 100。

# 设平均值为 x。
# 使 n 个仓库的库存数量相同 相当于

# 比 x 大的流出 a[i]−x ,
# 比 x 小的流入 x−a[i],
# 其余点 流入流出平衡。
# 故令比 x 大为源点 (限制流出 a[i]−x)，比 x 小的为汇点 (限制流入 x−a[i])。
# 其余为容量为inf流量为1的中间点(流入流出平衡)。
# 由于费用流是在最大流条件下的，于是以上条件必然得到满足。
# https://www.acwing.com/solution/content/26755/
# !加虚拟边平衡流量

# 如果不是环形 则为 https://leetcode.cn/problems/super-washing-machines/
from MinCostMaxFlow import MinCostMaxFlowDinic

INF = int(1e18)

n = int(input())  # n 个仓库 1≤n≤100,
stores = []  # 仓库的库存量
while len(stores) < n:
    stores.extend(map(int, input().split()))

avg = sum(stores) // n

START, END, OFFSET = 2 * n + 2, 2 * n + 3, n
mcmf = MinCostMaxFlowDinic(2 * n + 4, START, END)


for i, value in enumerate(stores):
    diff = value - avg
    if diff > 0:
        mcmf.addEdge(START, i, diff, 0)  # 流量守恒
    elif diff < 0:
        mcmf.addEdge(i, END, -diff, 0)

    mcmf.addEdge(i, (i + 1) % len(stores), INF, 1)  # 相邻边容量为inf，费用为1，表示运输单位货物需要1的费用。
    mcmf.addEdge(i, (i - 1) % len(stores), INF, 1)


print(mcmf.work()[1])
