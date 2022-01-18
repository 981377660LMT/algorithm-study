# 树，无环图
# 请问最多经过多少个单位时间后，小团会被追上。

# 答：
# 追的人远，被追的人近
# 只需要求出两个人到各个点的最短路径，追的那个人费时最长且比被追的人慢，则是答案

# 因为是无权图,bfs即可，不需要dijk(其实也可以)

from collections import defaultdict, deque
from heapq import heappop, heappush
from typing import Deque, List, Tuple

INF = 0x3F3F3F3F
n, x, y = list(map(int, input().split()))
adjMap = defaultdict(list)
for _ in range(n - 1):
    u, v = list(map(int, input().split()))
    adjMap[u].append(v)
    adjMap[v].append(u)


def getDistByBFS(start: int) -> List[int]:
    dist = [INF] * (n + 1)
    dist[start] = 0
    queue: Deque[Tuple[int, int]] = deque([(start, 0)])
    while queue:
        cur, steps = queue.popleft()
        for next in adjMap[cur]:
            if dist[cur] + 1 < dist[next]:
                dist[next] = dist[cur] + 1
                queue.append((next, steps + 1))
    return dist


def getDistByDijk(start: int) -> List[int]:
    dist = [INF] * (n + 1)
    dist[start] = 0
    queue = [(start, 0)]
    while queue:
        cur, steps = heappop(queue)
        for next in adjMap[cur]:
            if dist[cur] + 1 < dist[next]:
                dist[next] = dist[cur] + 1
                heappush(queue, (next, steps + 1))
    return dist


distX = getDistByBFS(x)
distY = getDistByBFS(y)
distX = getDistByDijk(x)
distY = getDistByDijk(y)
res = 0
for d1, d2 in zip(distX, distY):
    if d1 > d2:
        res = max(res, d1)

print(res)
