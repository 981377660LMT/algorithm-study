# 有 n 件工作要分配给 n 个人做。
# 第 i 个人做第 j 件工作产生的效益为 cij。
# 试设计一个将 n 件工作分配给 n 个人做的分配方案。
# 对于给定的 n 件工作和 n 个人，计算最优分配方案和最差分配方案。
# 1≤n≤50,

# !最优解肯定是最大匹配 即取最大流
# !最小费用最大流 最大费用最大流(边的费用取负)


from MinCostMaxFlow import MinCostMaxFlowDinic


n = int(input())
# 最小费用、最大费用
START, END, OFFSET = 2 * n + 2, 2 * n + 3, n
mcmf1 = MinCostMaxFlowDinic(2 * n + 5, START, END)
mcmf2 = MinCostMaxFlowDinic(2 * n + 5, START, END)

for i in range(n):
    nums = list(map(int, input().split()))
    for j, cost in enumerate(nums):
        mcmf1.addEdge(i, j + OFFSET, 1, cost)
        mcmf2.addEdge(i, j + OFFSET, 1, -cost)

for i in range(n):
    mcmf1.addEdge(START, i, 1, 0)
    mcmf2.addEdge(START, i, 1, 0)
    mcmf1.addEdge(i + OFFSET, END, 1, 0)
    mcmf2.addEdge(i + OFFSET, END, 1, 0)


_, cost1 = mcmf1.work()
_, cost2 = mcmf2.work()
print(cost1, -cost2, sep="\n")
