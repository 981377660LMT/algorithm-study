from collections import defaultdict, deque
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# bfs + 删除原来的边

if __name__ == "__main__":
    n, m = map(int, input().split())
    adjMap = defaultdict(dict)
    visited = [-1] * n
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u][v] = w
        adjMap[v][u] = w

    k = int(input())
    starts = list(map(int, input().split()))
    starts = [x - 1 for x in starts]
    for start in starts:
        visited[start] = 0
    queue = deque(starts)

    d = int(input())
    limits = list(map(int, input().split()))

    dep = 0
    while queue and dep < d:
        nextQueue = deque()
        len_ = len(queue)
        for _ in range(len_):
            cur = queue.popleft()
            for next in list(adjMap[cur]):
                w = adjMap[cur][next]
                if w <= limits[dep] and visited[next] == -1:
                    visited[next] = dep + 1
                    nextQueue.append(next)
                    del adjMap[next][cur]
                    del adjMap[cur][next]
        queue = nextQueue
        dep += 1

    for i in range(n):
        print(visited[i])
