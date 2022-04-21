# 无向图删去一个点能得到的连通块数量的最大值

# 给定一个由 n 个点 m 条边构成的无向图，请你求出该图删除一个点之后，连通块最多有多少。
from collections import defaultdict
from Tarjan import Tarjan

n, m = map(int, input().split())
adjMap = defaultdict(set)
for _ in range(m):
    u, v = map(int, input().split())
    adjMap[u].add(v)
    adjMap[v].add(u)

# https://www.acwing.com/solution/content/20702/
