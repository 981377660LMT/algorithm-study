from collections import deque
import sys
from typing import Deque, Set, Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 頂点に 1 から N の番号が、辺に 1 から M の番号がついた N 頂点 M 辺の単純無向グラフが与えられます。辺 i は頂点 u
# i
# ​
#   と頂点 v
# i
# ​
#   を結んでいます。また、各頂点の次数は 10 以下です。
# 頂点 1 を始点とする単純パス(同じ頂点を複数回通らないパス)の個数を K とします。min(K,10
# 6
#  ) を出力してください。

# bfs+记录路径
# !为什么要用dfs不能用bfs
if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u -= 1
        v -= 1
        adjList[u].append(v)
        adjList[v].append(u)

    res = 1
    queue: Deque[Tuple[int, Set[int]]] = deque()
    queue.append((0, {0}))
    while queue:
        cur, visted = queue.popleft()
        for next in adjList[cur]:
            if next in visted:
                continue
            res += 1
            if res >= 10**6:
                print(10**6)
                exit(0)
            queue.append((next, visted | {next}))
    print(res)
