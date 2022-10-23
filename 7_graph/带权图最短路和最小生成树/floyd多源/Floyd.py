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
