# F - Add One Edge 3
# https://atcoder.jp/contests/abc401/tasks/abc401_f
# 给定两棵树，树1和树2，将两棵树通过添加一条边合并后形成新树。求所有可能的连接方式下新树直径的总和。


from bisect import bisect_right
from itertools import accumulate
from typing import List, Tuple
from collections import deque

INF = int(1e18)


def getTreeDiameter(n: int, tree: List[List[int]], start=0) -> Tuple[int, List[int]]:
    """求无权树的(直径长度,直径路径)."""

    def dfs(start: int) -> Tuple[int, List[int]]:
        dist = [-1] * n
        dist[start] = 0
        stack = [start]
        while stack:
            cur = stack.pop()
            for next in tree[cur]:
                if dist[next] != -1:
                    continue
                dist[next] = dist[cur] + 1
                stack.append(next)
        endPoint = dist.index(max(dist))
        return endPoint, dist

    u, _ = dfs(start)
    v, dist = dfs(u)
    diameter = dist[v]
    path = [v]
    while u != v:
        for next in tree[v]:
            if dist[next] + 1 == dist[v]:
                path.append(next)
                v = next
                break

    return diameter, path


def bfs(start: int, adjList: List[List[int]]) -> List[int]:
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


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n1 = int(input())
    tree1 = [[] for _ in range(n1)]
    for _ in range(n1 - 1):
        u, v = map(int, input().split())
        tree1[u - 1].append(v - 1)
        tree1[v - 1].append(u - 1)
    n2 = int(input())
    tree2 = [[] for _ in range(n2)]
    for _ in range(n2 - 1):
        u, v = map(int, input().split())
        tree2[u - 1].append(v - 1)
        tree2[v - 1].append(u - 1)

    def f(tree: List[List[int]]) -> Tuple[int, List[int]]:
        d, path = getTreeDiameter(len(tree), tree)
        dist1 = bfs(path[0], tree)
        dist2 = bfs(path[-1], tree)
        res = [max(a, b) for a, b in zip(dist1, dist2)]
        return d, res

    d1, dist1 = f(tree1)
    d2, dist2 = f(tree2)
    if d1 < d2:
        d1, d2 = d2, d1
        dist1, dist2 = dist2, dist1

    # for v1 in dist1:
    #     for v2 in dist2:
    #         res += max(d2, v1 + v2 + 1)
    # 优化方法1：卷积

    # 优化方法2：排序(也可以滑窗，这里二分了)
    dist2.sort()
    presum2 = list(accumulate(dist2, initial=0))
    res = 0
    for v1 in dist1:
        threshold = d1 - v1 - 1
        pos = bisect_right(dist2, threshold)
        res += d1 * pos
        if pos < len(dist2):
            remainingSum = presum2[-1] - presum2[pos]
            remainingCount = len(dist2) - pos
            res += remainingSum + v1 * remainingCount + remainingCount
    print(res)
