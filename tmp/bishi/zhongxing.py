from collections import defaultdict
from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, k = map(int, input().split())
k -= 1
adjMap = defaultdict(lambda: defaultdict(lambda: INF))
for _ in range(n):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    adjMap[u][v] = min(adjMap[u][v], w)

pq = [(0, k)]
dist = [INF] * n
dist[k] = 0
while pq:
    curDist, cur = heappop(pq)
    if dist[cur] < curDist:
        continue
    for next in adjMap[cur]:
        if dist[next] > dist[cur] + adjMap[cur][next]:
            dist[next] = dist[cur] + adjMap[cur][next]
            heappush(pq, (dist[next], next))

max_ = max(dist)
print(max_ if max_ != INF else -1)
