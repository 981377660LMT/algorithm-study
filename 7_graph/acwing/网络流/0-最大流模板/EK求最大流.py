from collections import defaultdict
from Maxflow import MaxFlow

# 点的编号从 1 到 n。
# 输出点 S 到点 T 的最大流。
# 如果从点 S 无法到达点 T 则输出 0。
# 2≤n≤1000,
# 1≤m≤10000,

n, m, start, end = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(int))

# 从点 u 到点 v 存在一条有向边，容量为 c。
for _ in range(m):
    u, v, c = map(int, input().split())
    adjMap[u][v] = c

maxFlow = MaxFlow(adjMap)
print(maxFlow.calMaxFlow(start, end))
