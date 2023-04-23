from collections import deque
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    edges = []
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u -= 1
        v -= 1
        edges.append((u, v))
        adjList[u].append(v)
        adjList[v].append(u)

    k = int(input())
    limits = []
    for _ in range(k):
        p, d = map(int, input().split())
        p -= 1
        limits.append((p, d))

    def bfs(start: int, dep: int) -> None:
        queue = deque([start])
        visited = set([start])
        todo = dep
        while queue and todo > 0:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                mustWhite[cur] = True
                for next in adjList[cur]:
                    if next not in visited:
                        visited.add(next)
                        queue.append(next)
            todo -= 1

    def bfsDepth(start: int, dep: int) -> List[int]:
        queue = deque([start])
        visited = set([start])
        todo = dep
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

    # !一开始全部染黑
    mustWhite = [False] * n
    for p, d in limits:
        if d == 0:
            continue
        bfs(p, d)
    if all(mustWhite):
        print("No")
        exit(0)

    for p, d in limits:
        ps = bfsDepth(p, d)
        if all(mustWhite[ps[i]] for i in range(len(ps))):
            print("No")
            exit(0)
    print("Yes")
    res = [0 if mustWhite[i] else 1 for i in range(n)]
    print("".join(map(str, res)))
