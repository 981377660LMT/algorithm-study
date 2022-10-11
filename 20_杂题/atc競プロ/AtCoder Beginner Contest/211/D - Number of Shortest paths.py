# 无向图最短路计数

from collections import deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    count = [0] * n  # 记录到0达每个点的最短路条数
    dist = [INF] * n
    dist[0] = 0
    count[0] = 1
    queue = deque([(0, 0)])
    while queue:
        curDist, cur = queue.popleft()
        if curDist > dist[cur]:
            continue
        for next in adjList[cur]:
            cand = curDist + 1
            if dist[next] > cand:
                dist[next] = cand
                count[next] = count[cur]
                queue.append((cand, next))
            elif dist[next] == cand:
                count[next] += count[cur]
                count[next] %= MOD

    print(count[n - 1])
