"""
最小生成树带点权

要使0-n-1 连通
有三种建设方案
每个点以airCosti建设一个飞机场、每个点以minatoCosti建设一个港口、每个点以w连通u到v的道路
都具有 机场/港口/道路的点可以互相往来
求连通所有点的最小成本

1. slave 集群通信 => 引入 master - slave 解耦
2. 枚举4种情况 : 点群最后需不需要与 (n+1)/(n+2) 连通
   枚举不重不漏,看似不讲道理,却又无法反驳
   为什么这样就是对的? => 因为正好覆盖了所有情况
"""


from itertools import product
import sys
from 模板 import kruskal

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n, m = map(int, input().split())
    airport = list(map(int, input().split()))
    harbor = list(map(int, input().split()))
    road = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        road.append((u, v, w))

    FOO1, FOO2 = n + 1, n + 2  # 用于标记点
    res = INF
    for needAir, needHarbor in product([True, False], repeat=2):
        edges = road[:]
        for i in range(n):
            if needAir:
                edges.append((FOO1, i, airport[i]))
            if needHarbor:
                edges.append((FOO2, i, harbor[i]))
        cost, _ = kruskal(n + needAir + needHarbor, edges)
        if cost != -1:
            res = min(res, cost)
    print(res)
