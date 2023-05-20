# https://atcoder.jp/contests/agc018/tasks/agc018_d
# 起点任意，求访问树中所有点的最长距离

# 树哈密尔顿回路
# !以重心为根


from collections import deque
from heapq import nsmallest
from typing import List, Tuple

INF = int(1e18)


def getCenter(n: int, tree: List[List[Tuple[int, int]]], root=0) -> List[int]:
    """求重心."""

    def dfs(cur: int, pre: int) -> None:
        subsize[cur] = 1
        for next, _ in tree[cur]:
            if next == pre:
                continue
            dfs(next, cur)
            subsize[cur] += subsize[next]
            weight[cur] = max(weight[cur], subsize[next])
        weight[cur] = max(weight[cur], n - subsize[cur])
        if weight[cur] <= n / 2:
            res.append(cur)

    res = []
    weight = [0] * n  # 节点的`重量`，即以该节点为根的子树的最大节点数
    subsize = [0] * n  # 子树大小
    dfs(root, -1)
    return res


def bfs(tree: List[List[Tuple[int, int]]], start: int) -> List[int]:
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    while queue:
        cur = queue.popleft()
        for next, w in tree[cur]:
            if dist[next] == INF:
                dist[next] = dist[cur] + w
                queue.append(next)
    return dist


def maxDistTranverse(n: int, tree: List[List[Tuple[int, int]]]) -> int:
    """
    访问树中所有点的最长距离.
    tree是无向图邻接表.
    """
    center = getCenter(n, tree)
    root = center[0]
    dist = bfs(tree, root)
    if len(center) == 1:
        max2 = nsmallest(2, dist)
        return 2 * sum(dist) - sum(max2)
    return 2 * sum(dist) - dist[center[0]] - dist[center[1]]


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e7))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v, w = map(int, input().split())
        adjList[u - 1].append((v - 1, w))
        adjList[v - 1].append((u - 1, w))
    print(maxDistTranverse(n, adjList))
