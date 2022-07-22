# 有N个人和K个国家,每个人刚好属于一个国家:第i个人属于A。个国家此外,有L个名人在他们之中: B1, B2,. ..,B，
# 起初,没有任何两个人是朋友
# 对于这M 对人,三桥,作为上帝,能让他们通过一定花费成为朋友:对于每个1<i≤M ,他能支付C花费使得U,V成为朋友
# 现在,对1<i≤N ,解决如下问题
# 三桥君是否每个人都认识一个来自不同国家的名人如果可以,求最小花费
# !(求每一个人和他的外国名人做朋友（间接朋友也行）的最小花费。)


# 如果没有颜色限制:直接从名人开始多源最短路

# !1.我们把所有受欢迎的人放在优先队列里同时跑最短路，并记录他们所属的国家,
# 在这个过程中我们通过记录，保证每个国家只会更新aj一次。
# !⒉.那么当人a;被第一次和第二次取出的时候，所需的花费分别是最小，次小的，
# 并且更新这两次的国家不同。因为a;只会属于一个国家，所以这两个最小数中一定有一个和自己的国家不同。
# !3.如果最短路的类型是不同国籍的，就直接用最短；不然就切换到次短。

from collections import defaultdict
from heapq import heappop, heappush
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, m, k, l = map(int, input().split())
    adjMap = defaultdict(lambda: defaultdict(lambda: int(1e18)))
    country = [int(num) - 1 for num in input().split()]  # 每一个人属于哪个国家
    celeb = [int(num) - 1 for num in input().split()]  # k个有名的人的编号
    for _ in range(m):  # 再给m个关系，表示u和v做朋友的花费
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u][v] = min(adjMap[u][v], w)
        adjMap[v][u] = min(adjMap[v][u], w)

    dist = [[int(1e18), int(1e18)] for _ in range(n)]  # 最短路/次短路
    distCountry = [[-1, -1] for _ in range(n)]  # 最短路/次短路起点所属国家
    pq = []  # 距离 当前点 起点所属国家
    for start in celeb:
        pq.append((0, start, country[start]))
        dist[start][0] = 0
        distCountry[start][0] = country[start]

    while pq:
        cost, cur, root = heappop(pq)
        if dist[cur][1] < cost:  # 比次短路差
            continue

        for next, weight in adjMap[cur].items():
            cand = cost + weight
            if cand < dist[next][0]:
                dist[next][0] = cand
                distCountry[next][0] = root
                heappush(pq, (cand, next, root))
            elif cand < dist[next][1] and distCountry[next][0] != root:  # 最短次短路不能是同一个国家
                dist[next][1] = cand
                distCountry[next][1] = root
                heappush(pq, (cand, next, root))

    res = [-1] * n
    for i in range(n):
        dist0, country0 = dist[i][0], distCountry[i][0]
        if country0 != country[i]:
            res[i] = dist0 if dist0 < int(1e16) else -1
        else:
            res[i] = dist[i][1] if dist[i][1] < int(1e16) else -1
    print(*res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()


# TODO
# 哪里有问题
