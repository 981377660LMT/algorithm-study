# https://acexam.nowcoder.com/coding/?uid=20D4CBC7F99CAB36&qid=10520381
# 顺丰半决赛
# 地图上有个n城市，m条道路。每个道路连接着两个城市。顺丰小哥需要将快递从1号城市运到n号城市，
# 另外，顺丰小哥有一个魔法，他可以花费x时间从一个城市传送到任意另一个城市。
# !为了避免暴露自己会魔法的事实，他不会在起点和终点使用魔法（也不会作为魔法的目的地），且魔法最多只能使用一次。
# 顺丰小哥想知道，自己完成送快递至少需要多少时间？

# !传送可以等价于和虚拟结点相连,传送的边权变为w/2

from typing import Sequence, Tuple
from heapq import heappop, heappush

INF = int(1e18)


def dijkstraWithTeleportOnce(
    n: int, edges: Sequence[Tuple[int, int, int]], start: int, target: int, teleportCost: int
) -> int:
    """
    允许使用一次传送的最短路,传送代价为teleportCost.传送不能作为起点和终点,也不能作为传送的目的地.
    edges为无向边的列表,每个元素为(u, v, w)表示u和v之间有一条权重为w的边.
    """
    adjList = [[] for _ in range(n + 1)]
    for u, v, w in edges:
        w *= 2
        adjList[u].append((v, w))
        adjList[v].append((u, w))

    DUMMY = n
    for i in range(n):
        if i == start or i == target:
            continue
        adjList[i].append((DUMMY, teleportCost))
        adjList[DUMMY].append((i, teleportCost))

    dist = [[INF, INF] for _ in range(n + 1)]  # dist[i][j]表示从start到i,使用了j次传送的最短距离.
    dist[start][0] = 0
    pq = [(0, start, 0)]  # (dist, cur, count)
    while pq:
        curStep, cur, curCount = heappop(pq)
        if dist[cur][curCount] < curStep:
            continue
        if cur == target:
            return curStep // 2
        for next, weight in adjList[cur]:
            nextCount = curCount + (next == DUMMY)
            if nextCount > 1:
                continue
            nextDist = curStep + weight
            if nextDist < dist[next][nextCount]:
                dist[next][nextCount] = nextDist
                heappush(pq, (nextDist, next, nextCount))

    return -1


if __name__ == "__main__":
    n, m, x = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w))

    print(dijkstraWithTeleportOnce(n, edges, 0, n - 1, x))
