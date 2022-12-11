from itertools import product
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(1e18)
# 正整数 N,M が与えられます。
# 頂点に 1,…,N の番号が付けられた N 頂点の単純連結無向グラフであって、以下の条件を満たすものの総数を M で割った余りを求めてください。

# 全ての u=2,…,N−1 について、頂点 1 から頂点 u までの最短距離は、頂点 1 から頂点 N までの最短距離より真に小さい。
# ただし、頂点 u から頂点 v までの最短距離とは、頂点 u,v を結ぶ単純パスに含まれる辺の本数の最小値を指します。
# また、2 つのグラフが異なるとは、ある 2 頂点 u,v が存在して、これらの頂点を結ぶ辺が一方のグラフにのみ存在することを指します。


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
    from collections import deque

    def foo(n: int, mod: int) -> int:
        def check(adjMatrix: List[List[int]]) -> bool:
            adjList = [[] for _ in range(n)]
            for i in range(n):
                for j in range(i + 1, n):
                    if adjMatrix[i][j]:
                        adjList[i].append(j)
                        adjList[j].append(i)

            dist = bfs(0, adjList)
            return all(dist[i] < dist[n - 1] for i in range(n - 1)) and dist[n - 1] != INF

        res = 0
        for select in product(range(2), repeat=n * (n - 1) // 2):
            adjMatrix = [[0] * n for _ in range(n)]
            eid = 0
            for i in range(n):
                for j in range(i + 1, n):
                    adjMatrix[i][j] = select[eid]
                    adjMatrix[j][i] = select[eid]
                    eid += 1
            if check(adjMatrix):
                res = (res + 1) % mod
        return res

    for i in range(3, 10):
        print(foo(i, 998244353))
# 1
# 8
# 98
# 2296
# 107584
