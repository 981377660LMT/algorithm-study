# 水平的列车上有n个座位，从左到右座位号为1,2...n。
# !现在有m条规定，每条规定的形式如下∶从座位l到座位r，不多于x个人乘坐。
# 在满足所有规定的前提下，该列车最多能柔坐多少人?

# !最大值=>求最短路

from collections import defaultdict, deque
from typing import List, Mapping

INF = int(1e18)


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
    adjMap[v][u - 1] = max(adjMap[v][u - 1], -w)  # Su-1 - Sv >= -w

# !前缀和满足的约束
for i in range(1, n + 1):
    adjMap[i - 1][i] = max(adjMap[i - 1][i], 0)  # Si - Si-1 >= 0
    adjMap[i][i - 1] = max(adjMap[i][i - 1], -1)  # Si-1 - Si >= -1


dist = spfa(n + 1, adjMap, n)
print(-dist[0])

# 10 3
# 1 4 2
# 3 6 2
# 10 10 1
# 输出8
