from itertools import permutations
from math import ceil
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 最初
# N 頂点
# 0 辺の無向グラフがあり、各頂点には
# 1 から
# N まで番号がついています。
# Q 個のクエリが与えられるので、順に処理し、各クエリの後における「他のどの頂点とも辺で結ばれていない頂点」の数を出力してください。

# i 個目のクエリは
# query
# i
# ​
#   であり、各クエリは次の
# 2 種類いずれかです。

# 1 u v: 頂点
# u と頂点
# v を辺で結ぶ。このクエリが与えられる直前の時点で、頂点
# u と頂点
# v は辺で結ばれていない事が保証される。

# 2 v : 頂点
# v と他の頂点を結ぶ辺をすべて削除する。（頂点
# v 自体は削除しない。）
if __name__ == "__main__":
    n, q = map(int, input().split())
    res = n
    adjList = [set() for _ in range(n)]
    deg = [0] * n
    for _ in range(q):
        t, *args = map(int, input().split())
        if t == 1:
            u, v = args
            u, v = u - 1, v - 1
            adjList[u].add(v)
            adjList[v].add(u)
            if deg[u] == 0:
                res -= 1
            if deg[v] == 0:
                res -= 1
            deg[u] += 1
            deg[v] += 1
        else:
            v = args[0] - 1
            nexts = adjList[v]
            if deg[v] != 0:
                res += 1
            deg[v] = 0

            adjList[v] = set()  # 为什么clear不对???
            for next in nexts:
                deg[next] -= 1
                adjList[next].remove(v)
                if deg[next] == 0:
                    res += 1

        print(res)
