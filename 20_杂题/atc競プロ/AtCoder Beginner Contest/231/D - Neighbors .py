"""判断是否存在1-n的排列使得ai与bi相邻

所有点度数不超过2
不存在环
"""
from collections import deque
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def solve(n: int, edges: List[List[int]]) -> bool:
        """判断是否存在0-n-1的排列使得edges里的边相邻

        1. 所有点度数不超过2
        2. 无向图不存在环
        """
        adjList = [[] for _ in range(n)]
        deg = [0] * n
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
            deg[u] += 1
            deg[v] += 1

        if any(d > 2 for d in deg):
            return False

        queue = deque(i for i, d in enumerate(deg) if d == 1)
        while queue:
            cur = queue.popleft()
            for next in adjList[cur]:
                deg[next] -= 1
                if deg[next] == 1:
                    queue.append(next)

        return all(d == 0 for d in deg)

    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        a, b = map(int, input().split())
        a, b = a - 1, b - 1
        edges.append([a, b])
    print("Yes" if solve(n, edges) else "No")
