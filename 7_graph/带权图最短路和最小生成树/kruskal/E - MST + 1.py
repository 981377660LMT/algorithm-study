# 有一个无向连通图 M，每次询问给出一条边 e，询问 e 加入到 M 中是否会影响 M 的最小生成树。
# !把 M 中的边和查询边放到一起跑最小生成树算法，但是注意查询边只判断不改变连通性。


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, m, q = map(int, input().split())
edges1 = []
for _ in range(m):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    edges1.append((u, v, w))


edges2 = []
for _ in range(q):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    edges2.append((u, v, w))
