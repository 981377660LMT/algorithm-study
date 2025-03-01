# abc394-E - Palindromic Shortest Path-回文最短路
# https://atcoder.jp/contests/abc394/tasks/abc394_e
# 给定一个n个点的有向图，每条边上有一个字符，问从每个点到每个点的最短回文路径长度是多少。
# 如果不存在这样的路径，输出-1。
# n <= 100.
#
# !状态为 (i, j) 的无权图最短路.
# 初始时，队列为所有的 (i, i) + (i, j) 的状态，分别对应长度为偶数和奇数的回文串。
# 每次拓展会让回文串的长度+2, 复杂度 O(n^2).


from collections import defaultdict, deque

import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n = int(input())
    adjMatrix = [input() for _ in range(n)]

    graph = [defaultdict(list) for _ in range(n)]  # graph[u][c]
    revGraph = [defaultdict(list) for _ in range(n)]
    for x in range(n):
        for y, c in enumerate(adjMatrix[x]):
            if c != "-":
                graph[x][c].append(y)
                revGraph[y][c].append(x)

    dist = [[INF] * n for _ in range(n)]
    queue = deque()  # (from, to)

    for x in range(n):
        dist[x][x] = 0
        queue.append((x, x))

    for x in range(n):
        for ys in graph[x].values():
            for y in ys:
                if x != y:
                    dist[x][y] = 1
                    queue.append((x, y))

    while queue:
        x1, y1 = queue.popleft()
        d = dist[x1][y1]
        for c in graph[y1]:
            if c in revGraph[x1]:
                for x2 in revGraph[x1][c]:
                    for y2 in graph[y1][c]:
                        if dist[x2][y2] > d + 2:
                            dist[x2][y2] = d + 2
                            queue.append((x2, y2))

    res = []
    for vs in dist:
        cur = [str(v) if v != INF else "-1" for v in vs]
        res.append(" ".join(cur))
    print("\n".join(res))
