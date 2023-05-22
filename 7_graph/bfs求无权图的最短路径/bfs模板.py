# bfs模板 类似于dijkstra

from collections import deque
from typing import List, Sequence, Tuple

INF = int(1e18)


def bfs(start: int, adjList: List[List[int]]) -> List[int]:
    """时间复杂度O(V+E)"""
    n = len(adjList)
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            cand = dist[cur] + 1
            if cand < dist[next]:
                dist[next] = cand
                queue.append(next)
    return dist


# bfs求出一条路径
def bfs2(n: int, adjList: Sequence[Sequence[int]], start: int, end: int) -> Tuple[int, List[int]]:
    """bfs求出起点到end的(最短距离,路径) 时间复杂度O(V+E)"""
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    pre = [-1] * n  # 记录每个点的前驱

    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            cand = dist[cur] + 1
            if cand < dist[next]:
                dist[next] = cand
                pre[next] = cur
                queue.append(next)

    if dist[end] == INF:
        return INF, []

    path = []
    cur = end
    while pre[cur] != -1:
        path.append(cur)
        cur = pre[cur]
    path.append(start)
    return dist[end], path[::-1]


def bfsBanEdge(
    n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int, banEdge: int = -1
) -> Tuple[List[int], List[Tuple[int, int]]]:
    """bfs求最短路,并记录一条路径."""
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    pre = [(-1, -1)] * n  # (preNode, preEdge) bfs记录路径

    while queue:
        cur = queue.popleft()
        for next, eid in adjList[cur]:
            if eid == banEdge:
                continue
            cand = dist[cur] + 1
            if cand < dist[next]:
                pre[next] = (cur, eid)  # type: ignore
                dist[next] = cand
                queue.append(next)

    return dist, pre


from collections import deque
from typing import List


def bfsDepth(adjList: List[List[int]], start: int, dist: int) -> List[int]:
    """返回距离start为dist的结点"""
    if dist < 0:
        return []
    if dist == 0:
        return [start]
    queue = deque([start])
    visited = set([start])
    todo = dist
    while queue and todo > 0:
        len_ = len(queue)
        for _ in range(len_):
            cur = queue.popleft()
            for next in adjList[cur]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        todo -= 1
    return list(queue)


if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for i in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append((v, i))  # !记录每条边的编号

    dist, pre = bfsBanEdge(n, adjList, 0)
    if dist[n - 1] == INF:
        print(-1)
        exit(0)

    path = []  # !记录最短路上的边id
    cur = n - 1
    while cur != -1:
        path.append(pre[cur][1])
        cur = pre[cur][0]
    path.reverse()

    # E - Nearest Black Vertex
    # 给定一张图，要求给点涂黑白色，要求至少有一个黑点，且满足k个要求。
    # 每个 要求 (pi,di)表示点 pi距离黑点的最近距离恰好为 di。
    # 点数、边数 ≤2000
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    k = int(input())
    limits = []
    for _ in range(k):
        p, d = map(int, input().split())
        p -= 1
        limits.append((p, d))

    # !一开始全部染黑
    def fillWhite(start: int, dep: int) -> None:
        queue = deque([start])
        visited = set([start])
        while queue and dep:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                mustWhite[cur] = True
                for next in adjList[cur]:
                    if next not in visited:
                        visited.add(next)
                        queue.append(next)
            dep -= 1

    mustWhite = [False] * n
    for p, d in limits:
        if d == 0:
            continue
        fillWhite(p, d)
    if all(mustWhite):
        print("No")
        exit(0)

    for p, d in limits:
        level = bfsDepth(adjList, p, d)
        if all(mustWhite[level[i]] for i in range(len(level))):
            print("No")
            exit(0)
    print("Yes")
    res = [0 if mustWhite[i] else 1 for i in range(n)]
    print("".join(map(str, res)))
