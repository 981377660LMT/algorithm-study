# 7_graph/acwing/差分约束/362. 区间-前缀和+差分约束.py
# TLE TODO


# 前缀和+差分约束
import sys
from collections import defaultdict, deque
from typing import List, Mapping

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def spfa(n: int, adjMap: Mapping[int, Mapping[int, int]], start: int) -> List[int]:
    """spfa求单源最长路 图中无正环"""
    dist = [-INF] * n
    dist[start] = 0
    queue = deque([start])
    inQueue = [False] * n
    inQueue[start] = True

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False
        for next in adjMap[cur]:
            weight = adjMap[cur][next]
            cand = dist[cur] + weight
            if cand > dist[next]:  # 最长路
                dist[next] = cand
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] > dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)
    return dist


n, m = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(lambda: -INF))  # 最长路
for _ in range(m):
    u, v, w = map(int, input().split())  # 从1开始
    adjMap[u - 1][v] = max(adjMap[u - 1][v], w)  # Sv - Su-1 >= w

# !前缀和满足的约束
for i in range(1, n + 1):
    adjMap[i - 1][i] = max(adjMap[i - 1][i], 0)  # Si - Si-1 >= 0
    adjMap[i][i - 1] = max(adjMap[i][i - 1], -1)  # Si-1 - Si >= -1


dist = spfa(n + 1, adjMap, 0)
for pre, cur in zip(dist, dist[1:]):
    print(cur - pre, end=" ")
