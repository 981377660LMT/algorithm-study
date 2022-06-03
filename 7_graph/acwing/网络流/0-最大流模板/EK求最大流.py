from collections import defaultdict
import sys
from Maxflow import MaxFlow, Dinic, EK
import gc

gc.disable()

# 给定一个包含 n 个点 m 条边的有向图，并给定每条边的容量，边的容量非负。
# 图中可能存在重边和自环。求从点 S 到点 T 的最大流。

# 点的编号从 1 到 n。
# 输出点 S 到点 T 的最大流。
# 如果从点 S 无法到达点 T 则输出 0。

# !EK:
# 2≤n≤1e3,
# 1≤m≤1e4,
# !DINIC:
# 2≤n≤1e4,
# 1≤m≤1e5,


# 图中可能存在重边和自环
input = sys.stdin.readline
n, m, start, end = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(int))

# 从点 u 到点 v 存在一条有向边，容量为 c。
for _ in range(m):
    u, v, c = map(int, input().split())
    adjMap[u][v] += c  # 可能存在重边

maxFlow = MaxFlow(adjMap)
maxFlow.switchTo(EK(adjMap))
print(maxFlow.calMaxFlow(start, end))
maxFlow.switchTo(Dinic(adjMap))
print(maxFlow.calMaxFlow(start, end))

