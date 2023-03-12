# https://www.acwing.com/activity/content/problem/content/920/

# 请你求出 1 号点到 n 号点的最短距离，如果无法从 1 号点走到 n 号点，则输出 impossible。
# 边权可能为负数。

# spfa可以过很多dijk的题
# 但是网格的图容易卡spfa
# 有边数限制也能用spfa，spfa本质就是让bf不去枚举到不可能会拓展的边，bf能做的spfa都能

# 1. dist
# 2. queue
# 3. 判重数组 inQueue

from collections import deque
from typing import List, Optional, Tuple

INF = int(1e18)


def spfa1(n: int, adjList: List[List[Tuple[int, int]]], start: int) -> Optional[List[int]]:
    """spfa求单源最短路(图中有负边无负环,如果有负环则返回None)

    适用于解决带有负权重的图,是Bellman-ford的常数优化版
    """
    dist = [INF] * n
    queue = deque([start])
    inQueue = [False] * n
    count = [0] * n

    inQueue[start] = True
    dist[start] = 0
    count[start] = 1

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:  # 如果要最长路这里需要改成 >
                dist[next] = cand
                if not inQueue[next]:
                    count[next] += 1
                    if count[next] >= n:
                        return
                    inQueue[next] = True
                    if queue and dist[next] < dist[queue[0]]:  # !酸辣粉优化 如果要最长路这里需要改成 >
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return dist


if __name__ == "__main__":
    #  spfa 求最短路
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v, w = list(map(int, input().split()))
        u, v = u - 1, v - 1
        adjList[u].append((v, w))
    dist = spfa1(n, adjList, 0)
    if dist is None:
        print("impossible")
        exit(0)
    print(dist[n - 1] if dist[n - 1] != INF else "impossible")
