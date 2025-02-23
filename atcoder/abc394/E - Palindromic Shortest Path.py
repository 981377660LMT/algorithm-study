import sys
from heapq import heappop, heappush


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n = int(input())
    adjMatrix = [input() for _ in range(n)]

    graph = [[] for _ in range(n)]
    revGraph = [[] for _ in range(n)]
    for u in range(n):
        for j in range(n):
            c = adjMatrix[u][j]
            if c != "-":
                graph[u].append((j, c))
                revGraph[j].append((u, c))

    dp = [[INF] * n for _ in range(n)]
    pq = []  # (dist, from, to)

    for u in range(n):
        dp[u][u] = 0
        heappush(pq, (0, u, u))
    for u in range(n):
        for v, c in graph[u]:
            if dp[u][v] > 1:
                dp[u][v] = 1
                heappush(pq, (1, u, v))

    while pq:
        d, x, y = heappop(pq)
        if d > dp[x][y]:
            continue
        for u, c1 in revGraph[x]:
            for v, c2 in graph[y]:
                if c1 == c2:
                    nd = d + 2
                    if nd < dp[u][v]:
                        dp[u][v] = nd
                        heappush(pq, (nd, u, v))

    res = []
    for u in range(n):
        cur = []
        for j in range(n):
            if dp[u][j] < INF:
                cur.append(str(dp[u][j]))
            else:
                cur.append(str(-1))
        res.append(" ".join(cur))
    print("\n".join(res))
