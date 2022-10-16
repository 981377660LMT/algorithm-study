# 每个判断为 right-left=dist
# 问所有的判断是否无矛盾

# !dfs或bfs求出每个点到每个组的根的距离,再逐一检验

from collections import deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    edges = []
    for _ in range(m):
        left, right, weigth = map(int, input().split())
        left, right = left - 1, right - 1
        adjList[left].append((right, weigth))
        adjList[right].append((left, -weigth))
        edges.append((left, right, weigth))

    # def dfs(cur: int, curDist: int) -> None:
    #     if visited[cur]:
    #         return
    #     visited[cur] = True
    #     dist[cur] = curDist
    #     for next, weight in adjList[cur]:
    #         dfs(next, curDist + weight)

    def bfs(start: int) -> None:
        queue = deque([start])
        visited[start] = True
        while queue:
            cur = queue.popleft()
            for next, weight in adjList[cur]:
                if not visited[next]:
                    visited[next] = True
                    dist[next] = dist[cur] + weight
                    queue.append(next)

    visited = [False] * n
    dist = [0] * n
    for i in range(n):
        if not visited[i]:
            # dfs(i, 0)
            bfs(i)

    for u, v, w in edges:
        if dist[u] + w != dist[v]:
            print("No")
            exit(0)

    print("Yes")
