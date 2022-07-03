# 请你求出 1 号点到 n 号点的最短距离，如果无法从 1 号点走到 n 号点，则输出 impossible。
# 边权可能为负数。

# spfa可以过很多dijk的题
# 但是网格的图容易卡spfa
# 有边数限制也能用spfa，spfa本质就是让bf不去枚举到不可能会拓展的边，bf能做的spfa都能

# 1. dist
# 2. queue
# 3. 判重数组 inQueue

from collections import defaultdict, deque
from typing import DefaultDict, Hashable, TypeVar

T = TypeVar("T", bound=Hashable)


def spfa(adjMap: DefaultDict[T, DefaultDict[T, int]], start: T) -> DefaultDict[T, int]:
    """spfa求单源最短路,适用于解决带有负权重的图,是Bellman-ford的常数优化版"""
    dist = defaultdict(lambda: int(1e18), {start: 0})
    queue = deque([start])
    inQueue = defaultdict(lambda: False, {start: True})

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False

        for next, weight in adjMap[cur].items():
            if dist[cur] + weight < dist[next]:
                dist[next] = dist[cur] + weight
                if not inQueue[next]:
                    inQueue[next] = True
                    # !酸辣粉优化
                    if queue and dist[next] < dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return dist


# !酸辣粉优化
# https://blog.csdn.net/zhouchangyu1221/article/details/90549195
# SLF(Small Label First) 优化
# 将原队列改成双端队列，对要加入队列的点 p，如果 dist[p] 小于队头元素 u 的 dist[u]，将其插入到队头，否则插入到队尾。
# SLF(Small Label First) 双端队列优化，也被戏称为“酸辣粉优化”
# SLF优化（酸辣粉优化），可以一定程度上（约20%）改善spfa的运行速度
