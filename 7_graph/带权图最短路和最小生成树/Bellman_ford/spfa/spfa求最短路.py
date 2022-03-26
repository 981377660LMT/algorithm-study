# 请你求出 1 号点到 n 号点的最短距离，如果无法从 1 号点走到 n 号点，则输出 impossible。
# 边权可能为负数。

# spfa可以过很多dijk的题
# 但是网格的图容易卡spfa
# 有边数限制也能用spfa，spfa本质就是让bf不去枚举到不可能会拓展的边，bf能做的spfa都能

from collections import defaultdict, deque
from typing import DefaultDict, List


n, m = map(int, input().split())
adjMap = defaultdict(list)
for _ in range(m):
    u, v, w = list(map(int, input().split()))
    adjMap[u].append((v, w))

# 1. dist
# 2. 队列
# 3. 判重数组


def spfa(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int) -> List[int]:
    """spfa求单源最短路，适用于解决带有负权重的图，是Bellman-ford的常数优化版"""
    dist = [int(1e20)] * (n + 1)
    dist[start] = 0

    queue = deque([start])
    isInqueue = [False] * (n + 1)  # 在队列里的点
    isInqueue[start] = True

    while queue:
        cur = queue.popleft()
        isInqueue[cur] = False  # 点从队列出来了

        # 更新过谁，就拿谁去更新别人
        for next, weight in adjMap[cur].items():
            if dist[cur] + weight < dist[next]:
                dist[next] = dist[cur] + weight
                if not isInqueue[next]:
                    isInqueue[next] = True
                    # 队列不为空，且当前元素距离小于队头，则加入队头，否则加入队尾
                    if queue and dist[next] < dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return dist


# dist = spfa(n, adjMap, 1)
# print(dist if dist[0] < int(1e20) else 'impossible')
