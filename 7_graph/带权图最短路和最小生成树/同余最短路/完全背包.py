# 给定 n 个整数，求这 n 个整数能拼凑出多少的其他整数（n 个整数可以重复取）
# 有n个物品,第i种物品单个体积为vi,价值为ci
# q次询问,每次给出背包的容积,你需要选择若干个物品,
# 每种物品可以选任意多个
# !在选出物品的体积的和恰好为 V 的前提下最大化选出物品的价值的和。
# 若不存在体积和恰好为 V 的方案，输出 -1
# 为了体现你解决 NP-Hard 问题的能力，V 会远大于 vi，详见数据范围部分。

# n<=50 vi<=1e5 ci<=1e5 q<=1e5
# !1e11<=V<=1e12

# 同余最长路
# 用性价比最高的物品的体积作为base,然后求多出来的价值的`同余最长路`
# !最后缺少的体积全用性价比最高的物品来补

from functools import cmp_to_key
from heapq import heappop, heappush

INF = int(1e18)

n, q = map(int, input().split())
goods = list(tuple(map(int, input().split())) for _ in range(n))  # !vi,ci
goods.sort(key=cmp_to_key(lambda x, y: x[0] * y[1] - x[1] * y[0]))  # 性价比降序

# !用性价比最高的物品的体积作为base,然后求价值的`同余最长路`
baseV, baseC = goods[0]
edges = []
for v, c in goods[1:]:
    edges.append((v % baseV, c - v // baseV * baseC))  # 取这个物品带来的额外收益


# dijk
dist = [-INF] * baseV
dist[0] = 0
pq = [(0, 0)]

while pq:
    curVal, cur = heappop(pq)
    curVal = -curVal
    if dist[cur] > curVal:
        continue
    for v, c in edges:
        next = (cur + v) % baseV
        cand = (curVal + c) - (cur + v) // baseV * baseC
        if cand > dist[next]:
            dist[next] = cand
            heappush(pq, (-cand, next))

# !最后缺少的体积全用性价比最高的物品来补
for _ in range(q):
    V = int(input())
    if dist[V % baseV] == -INF:
        print(-1)
    else:
        print(dist[V % baseV] + (V // baseV) * baseC)
