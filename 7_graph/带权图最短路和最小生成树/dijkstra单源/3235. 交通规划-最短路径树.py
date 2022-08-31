# 将现有的一部分铁路改造成高速铁路，使得任何两个城市间都可以通过高速铁路到达，
# 而且从所有城市乘坐高速铁路到首都的最短路程和原来一样长。
# 在这些条件下最少要改造多长的铁路。
# 所有的城市由 1 到 n 编号，首都为 1 号。
# 1≤n≤1e4 ，1≤m≤1e5


# 翻译：
# !给我们一个无向图，我们需要从无向图中选出来若干条边，使得选出来边的总权值和最小，
# !并且满足所有点到1号点的最短距离不变。
# 注意不是最小生成树 是最短路径树SPT
# !最短路径树SPT（Short Path Tree）是网络的源点到所有结点的最短路径构成的树。

# 总结：
# 1.先求一下1号点到所有点的最短距离（单源最短路算法）
# 2.然后再求一下每个点有哪些边满足条件`d[a] = d[b] + w(b->a)`（找备选边）
# 3.每个点所有满足条件的边中取一个权值最小的累加起来即可。

from collections import defaultdict
from heapq import heappop, heappush

INF = int(4e18)
n, m = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
for _ in range(m):
    u, v, w = map(int, input().split())
    adjMap[u][v] = min(adjMap[u][v], w)  # !注意重边
    adjMap[v][u] = min(adjMap[v][u], w)

dist = defaultdict(lambda: INF, {1: 0})
pq = [(0, 1)]
while pq:
    curDist, cur = heappop(pq)
    if dist[cur] < curDist:
        continue
    for next in adjMap[cur]:
        if dist[next] > dist[cur] + adjMap[cur][next]:
            dist[next] = dist[cur] + adjMap[cur][next]
            heappush(pq, (dist[next], next))

res = 0
for next in adjMap:
    cand = INF
    for cur, weight in adjMap[next].items():
        # 如果当前点k到1的距离等于点j到1的距离+点j到k的距离w，说明jk可为高速公路
        # !注意这里要反向推 正向推不对
        if dist[next] == dist[cur] + weight:
            cand = min(cand, weight)
    res += cand if cand != INF else 0

print(res)

# 2 4
# 3 5
# 4 2
