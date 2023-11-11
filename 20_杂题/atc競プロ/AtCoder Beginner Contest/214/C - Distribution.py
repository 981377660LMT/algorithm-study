"""宝石传递
环上传递宝石 求每个人在什么时候拿到第一个宝石

https://atcoder.jp/contests/abc214/editorial/2451
解法1:建图求虚拟点到每个点的最短路

解法2:dp递推
dp0=T0 dpi+1=min(dpi+Si,Ti+1)
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    # n = int(input())
    # cost = list(map(int, input().split()))  # !宝石传递到下一个人的时间(i=>(i+1)%n)))
    # getTime = list(map(int, input().split()))  # !给每个人发宝石的时间
    # min_ = min(getTime)
    # first = getTime.index(min_)
    # dp = [INF] * n
    # for i in range(n):
    #     pre, cur = (i + first - 1) % n, (i + first) % n
    #     dp[cur] = min(dp[pre] + cost[pre], getTime[cur])
    # print(*dp, sep="\n")

    #############################################################################
    n = int(input())
    cost = list(map(int, input().split()))  # !宝石传递到下一个人的时间(i=>(i+1)%n)))
    getTime = list(map(int, input().split()))  # !给每个人发宝石的时间
    from typing import List, Mapping
    from collections import defaultdict
    from heapq import heappop, heappush

    def dijkstra(n: int, adjMap: Mapping[int, Mapping[int, int]], start: int) -> List[int]:
        dist = [INF] * n
        dist[start] = 0
        pq = [(0, start)]

        while pq:
            curDist, cur = heappop(pq)
            if dist[cur] < curDist:
                continue
            for next in adjMap[cur]:
                cand = dist[cur] + adjMap[cur][next]
                if cand < dist[next]:
                    dist[next] = cand
                    heappush(pq, (dist[next], next))
        return dist

    adjMap = defaultdict(lambda: defaultdict(lambda: INF))
    for i in range(n):
        adjMap[i][(i + 1) % n] = cost[i]
        adjMap[n][i] = getTime[i]  # n为虚拟点
    print(*dijkstra(n + 1, adjMap, n)[:-1], sep="\n")
