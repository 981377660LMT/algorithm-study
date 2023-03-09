"""
1. 有些场合可以搭配bitset(状态压缩)达到64倍高速化效果
   floyd算法中的第一层循环的含义就是仅考虑1..k号点时,各个点之间的距离,
   这个距离可以换成连通性,这样可以用 bitset优化运算
    连通性求法1:O(V*E) bfs求最短路
    连通性求法2:O(V^3/64) 传递闭包(见下面的4)
"""

from itertools import product
from typing import List


INF = int(1e20)


if __name__ == "__main__":
    n = 3
    edges = [[0, 1, 2], [1, 0, 3], [2, 1, 1]]  # 无向带权边

    # !1.Floyd 求多源最短路
    dist = [[INF] * n for _ in range(n)]
    for i in range(n):
        dist[i][i] = 0

    for u, v, w in edges:
        dist[u][v] = dist[v][u] = w

    # !dis[k][i][j] 表示「经过若干个编号不超过 k 的节点」时，从 i 到 j 的最短路长度
    for k in range(n):  # !经过的中转点
        for i in range(n):
            for j in range(n):
                # !松弛：如果一条边可以被松弛了，说明这条边就没有必要留下了
                cand = dist[i][k] + dist[k][j]
                dist[i][j] = cand if dist[i][j] > cand else dist[i][j]

    # !2.Floyd 最短路计数 (Floyd求两点之间的最短路数)
    dist = [[INF] * n for _ in range(n)]
    for i in range(n):
        dist[i][i] = 0

    count = [[0] * n for _ in range(n)]
    for u, v, w in edges:
        dist[u][v] = dist[v][u] = w
        count[u][v] = count[v][u] = 1
    for k, i, j in product(range(n), repeat=3):
        cand = dist[i][k] + dist[k][j]
        if dist[i][j] == cand:
            count[i][j] += count[i][k] * count[k][j]
        elif dist[i][j] > cand:
            dist[i][j] = cand
            count[i][j] = count[i][k] * count[k][j]

    # !3.Floyd 记录最短路路径
    dist = [[INF] * n for _ in range(n)]
    for i in range(n):
        dist[i][i] = 0

    pre = [[i] * n for i in range(n)]  # pre[a][b]表示a作为最短路起点,取得最短路时b的前驱节点
    for u, v, w in edges:
        dist[u][v] = dist[v][u] = w
    for k, i, j in product(range(n), repeat=3):
        cand = dist[i][k] + dist[k][j]
        if dist[i][j] > cand:
            dist[i][j] = cand
            pre[i][j] = pre[k][j]

    def getPath(start: int, target: int) -> List[int]:
        cur = target
        res = [target]
        while cur != start:
            cur = pre[start][cur]
            res.append(cur)
        return res[::-1]

    # !4.Floyd + bitset 传递闭包维护有向图连通性 (第k轮连通时距离就为k)
    dp = [1 << i for i in range(n)]  # 每个点出发可以到达的点
    for u, v in edges:
        dp[u] |= 1 << v
    for k in range(n):
        for i in range(n):
            if dp[i] & (1 << k):
                dp[i] |= dp[k]
        for i in range(q):  # 查询距离
            ...
