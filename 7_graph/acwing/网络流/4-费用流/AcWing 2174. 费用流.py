# 给定一个包含 n 个点 m 条边的有向图，并给定每条边的容量和费用，边的容量非负。
# 图中可能存在重边和自环，保证费用不会存在负环。
# !求从 S 到 T 的最大流，以及在流量最大时的最小费用。
# 2≤n≤5000,
# 1≤m≤50000,
from collections import defaultdict
from MinCostMaxFlow import MinCostMaxFlow

if __name__ == '__main__':
    import sys

    # !图中存在重边和自环
    input = sys.stdin.readline
    n, m, start, end = map(int, input().split())
    mcmf = MinCostMaxFlow(start, end)

    # 从点 u 到点 v 存在一条有向边，容量为 c。
    for _ in range(m):
        u, v, c, cost = map(int, input().split())
        mcmf.addEdge(u, v, c, cost)
    flow, cost = mcmf.work()
    print(flow, cost, sep=' ')

