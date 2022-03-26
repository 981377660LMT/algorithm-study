# https://blog.csdn.net/zhouchangyu1221/article/details/90549195
# SLF(Small Label First) 优化
# 将原队列改成双端队列，对要加入队列的点 p，如果 dist[p] 小于队头元素 u 的 dist[u]，将其插入到队头，否则插入到队尾。
from collections import defaultdict, deque

# SLF(Small Label First) 双端队列优化，也被戏称为“酸辣粉优化”
# SLF优化（酸辣粉优化），可以一定程度上（约20%）改善spfa的运行速度


def spfa(n: int, adjMap: defaultdict, start: int, target: int) -> int:
    """spfa求单源最短路，适用于解决带有负权重的图，是Bellman-ford的常数优化版"""
    dist = [int(1e20)] * (n)
    dist[start] = 0

    queue = deque([start])
    isInqueue = [False] * (n)
    isInqueue[start] = True

    while queue:
        cur = queue.popleft()
        isInqueue[cur] = False

        for next, weight in adjMap[cur]:
            if dist[cur] + weight < dist[next]:
                dist[next] = dist[cur] + weight
                if not isInqueue[next]:
                    isInqueue[next] = True
                    if queue and dist[next] < dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return dist[target]
