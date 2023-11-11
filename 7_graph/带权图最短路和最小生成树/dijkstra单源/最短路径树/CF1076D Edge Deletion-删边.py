# 给定一个 n个点m边的无向连通带权图,要求删边至最多剩余 k 条边。
# 定义好点是指删边后, 1号节点到它的最短路长度仍然等于原图最短路长度的节点。
# 你需要最大化删边后的好点个数。
# 输出需要保留的边数和边的编号
# n,m<=3e5

# 好点的定义就是最短路径树上的点
# 需要保留的就是最短路径树上的边
# 如果最短路上的边数小于等于 k，那么我们直接把最短路树上的边全部保留，
# 否则我们在最短路径树上选择 k 条边留下
# 由于需要保证删掉以后图是联通的，所以我们从根节点 开始遍历，找到 k 条边留下即可。

from collections import deque
from heapq import heappop, heappush
from typing import List, Sequence, Tuple

INF = int(1e18)


def solve(n: int, edges: List[Tuple[int, int, int]], k: int, root=0) -> List[int]:
    adjList = [[] for _ in range(n)]
    for i, (u, v, w) in enumerate(edges):
        adjList[u].append((v, w, i))
        adjList[v].append((u, w, i))
    *_, preE = dijkstraSPT(n, adjList, root)

    # 从根节点遍历找k条在最短路径树上的边
    def bfs(cur: int, pre: int) -> None:
        queue = deque([(cur, pre)])
        while queue:
            cur, pre = queue.popleft()
            for next, _, eid in adjList[cur]:
                if next == pre:
                    continue
                if okEdge[eid] and len(res) < k:
                    res.append(eid)
                    queue.append((next, cur))

    okEdge = [False] * len(edges)
    for e in preE:
        if e != -1:
            okEdge[e] = True
    res = []
    bfs(root, -1)
    return res


# !dijkstra求出路径上每个点的前驱点和前驱边
# 其中邻接表的每个元素是一个三元组，分别是邻接点，边权，边的编号
def dijkstraSPT(
    n: int, adjList: Sequence[Sequence[Tuple[int, int, int]]], start: int
) -> Tuple[List[int], List[int], List[int]]:
    dist = [INF] * n
    preV, preE = [-1] * n, [-1] * n
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight, eid in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                preE[next] = eid
                preV[next] = cur
                heappush(pq, (dist[next], next))
            elif cand == dist[next]:
                preE[next] = eid
                preV[next] = cur
    return dist, preV, preE


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m, k = map(int, input().split())
    edges = []
    for i in range(m):
        u, v, w = map(int, input().split())
        edges.append((u - 1, v - 1, w))
    eids = solve(n, edges, k)
    print(len(eids))
    print(*[eid + 1 for eid in eids])
