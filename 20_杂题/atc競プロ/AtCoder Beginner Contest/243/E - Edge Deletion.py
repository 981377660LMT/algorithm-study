# n<=300 暗示floyd
# !有多少条边可以删除后不影响任意两点之间的最短距离和图的连通性(即不让任意两个结点之间的最短路变长.)
# !即可以删掉被松弛的边

# 首先如果某条边连接了a和b两个顶点如果存在另外一条a到b的路径距离小于等于这条边的边权,
# 那么这条边肯是可以删掉的,因为可以用这条路径来代替.
# !Floyd 最短路计数 (Floyd求两点之间的最短路数)


# 注意:
# !1. 不用dict 用数组存图/距离会快很多
# !2. 不用product 会快很多

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, m = map(int, input().split())

    edges = []
    dist = [[int(1e20)] * n for _ in range(n)]
    for i in range(n):
        dist[i][i] = 0

    count = [[0] * n for _ in range(n)]  # !最短路计数
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        dist[u][v] = dist[v][u] = w
        count[u][v] = count[v][u] = 1
        edges.append((u, v, w))

    # pypy3不要用product 会慢很多
    for k in range(n):
        for i in range(n):
            for j in range(n):
                cand = dist[i][k] + dist[k][j]
                if dist[i][j] == cand:
                    count[i][j] += count[i][k] * count[k][j]
                elif dist[i][j] > cand:
                    dist[i][j] = cand
                    count[i][j] = count[i][k] * count[k][j]

    res = 0
    for u, v, w in edges:
        # !被松弛的边或者没被松弛但最短路不止一条的边 可以删除
        if dist[u][v] < w or count[u][v] > 1:
            res += 1
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
