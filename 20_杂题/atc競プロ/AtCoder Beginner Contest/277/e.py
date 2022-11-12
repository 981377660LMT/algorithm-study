from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 個の頂点と M 本の辺からなる無向グラフが与えられます。
# i=1,2,…,M について、i 番目の辺は頂点 u
# i
# ​
#   と頂点 v
# i
# ​
#   を結ぶ無向辺であり、 a
# i
# ​
#  =1 ならばはじめは通行可能、a
# i
# ​
#  =0 ならばはじめは通行不能です。 また、頂点 s
# 1
# ​
#   、頂点 s
# 2
# ​
#   、… 、頂点 s
# K
# ​
#   の K 個の頂点にはスイッチがあります。

# 高橋君は、はじめ頂点 1 におり、「下記の移動とスイッチを押すの 2 つの行動のどちらかを行うこと」を好きなだけ繰り返します。

# 移動 : いまいる頂点と通行可能な辺を介して隣接する頂点を 1 つ選び、その頂点に移動する。
# スイッチを押す : いまいる頂点にスイッチがあるならば、そのスイッチを押す。その結果、グラフ上のすべての辺の通行可能・通行不能の状態が反転する。すなわち、通行可能である辺は通行不能に、通行不能である辺は通行可能に変化する。
# 高橋君が頂点 N に到達することが可能かどうかを判定し、可能な場合は頂点 N に到達するまでに行う移動の回数としてあり得る最小値を出力してください。
if __name__ == "__main__":
    n, m, k = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append((v, w))
        adjList[v].append((u, w))
    switches = set(map(lambda x: int(x) - 1, input().split()))

    # heap
    dist = [[INF, INF] for _ in range(n)]
    dist[0][0] = 0
    pq = [(0, 0, 0)]  # (移动次数，当前位置，当前flip)
    while pq:
        curDist, curPos, curState = heappop(pq)
        if curDist > dist[curPos][curState]:
            continue

        # 翻转
        if curPos in switches:
            nextState = 1 - curState
            if dist[curPos][nextState] > curDist:
                dist[curPos][nextState] = curDist
                heappush(pq, (curDist, curPos, nextState))

        # 移动
        for next, w in adjList[curPos]:
            if w ^ curState == 0:
                continue
            cand = curDist + 1
            if dist[next][curState] > cand:
                dist[next][curState] = cand
                heappush(pq, (cand, next, curState))

    res = min(dist[n - 1])
    print(res if res != INF else -1)
