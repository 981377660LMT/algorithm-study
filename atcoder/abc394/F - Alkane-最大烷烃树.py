from collections import deque
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    adjList = [[] for _ in range(n)]
    deg = [0] * n
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)
        deg[u] += 1
        deg[v] += 1

    isD4 = [d >= 4 for d in deg]
    maxD4 = 0
    visited = [False] * n
    for i in range(n):
        if isD4[i] and not visited[i]:
            count = 0
            queue = deque([i])
            visited[i] = True
            while queue:
                cur = queue.popleft()
                count += 1
                for next_ in adjList[cur]:
                    if isD4[next_] and not visited[next_]:
                        visited[next_] = True
                        queue.append(next_)
            maxD4 = max(maxD4, count)

    if maxD4 == 0:
        print(-1)
    else:
        print(3 * maxD4 + 2)
