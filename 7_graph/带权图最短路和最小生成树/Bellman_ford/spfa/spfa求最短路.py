# https://www.acwing.com/activity/content/problem/content/920/

# 请你求出 1 号点到 n 号点的最短距离，如果无法从 1 号点走到 n 号点，则输出 impossible。
# 边权可能为负数。

# spfa可以过很多dijk的题
# 但是网格的图容易卡spfa
# 有边数限制也能用spfa，spfa本质就是让bf不去枚举到不可能会拓展的边，bf能做的spfa都能

# 1. dist
# 2. queue
# 3. 判重数组 inQueue

from collections import defaultdict, deque
from typing import List, Mapping, Tuple

INF = int(1e18)


def spfa1(n: int, adjMap: Mapping[int, Mapping[int, int]], start: int) -> List[int]:
    """spfa求单源最短路(图中有负边无负环)

    适用于解决带有负权重的图,是Bellman-ford的常数优化版
    """
    dist = [INF] * n
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
            if cand < dist[next]:  # 如果要最长路这里需要改成 >
                dist[next] = cand
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] < dist[queue[0]]:  # !酸辣粉优化 如果要最长路这里需要改成 >
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return dist


def spfa(n: int, adjMap: Mapping[int, Mapping[int, int]]) -> Tuple[bool, List[int]]:
    """spfa求虚拟节点为起点的单源最长路 并检测正环"""
    dist = [0] * n
    queue = deque(list(range(n)))
    count = [0] * n
    inQueue = [True] * n

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False
        for next in adjMap[cur]:
            weight = adjMap[cur][next]
            cand = dist[cur] + weight
            if cand > dist[next]:
                dist[next] = cand
                count[next] = count[cur] + 1
                if count[next] >= n:
                    return False, []
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] > dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return True, dist


# !酸辣粉优化
# https://blog.csdn.net/zhouchangyu1221/article/details/90549195
# SLF(Small Label First) 优化
# 将原队列改成双端队列，对要加入队列的点 p，如果 dist[p] 小于队头元素 u 的 dist[u]，将其插入到队头，否则插入到队尾。
# SLF(Small Label First) 双端队列优化，也被戏称为“酸辣粉优化”
# SLF优化（酸辣粉优化），可以一定程度上（约20%）改善spfa的运行速度

if __name__ == "__main__":
    #  spfa 求最短路
    n, m = map(int, input().split())
    adjMap = defaultdict(lambda: defaultdict(lambda: INF))
    for _ in range(m):
        u, v, w = list(map(int, input().split()))
        u, v = u - 1, v - 1
        adjMap[u][v] = min(adjMap[u][v], w)
    dist = spfa1(n, adjMap, 0)
    print(dist[n - 1] if dist[n - 1] != INF else "impossible")
