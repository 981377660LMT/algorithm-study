# 地图上有个城市，条道路。每个道路连接着两个城市。顺丰小哥需要将快递从1号城市运到号城市，
# 另外，顺丰小哥有一个魔法，他可以花费时间从一个城市传送到任意另一个城市。为了避免暴露自己会魔法的事实，他不会在起点和终点使用魔法（也不会作为魔法的目的地），且魔法最多只能使用一次。
# 顺丰小哥想知道，自己完成送快递至少需要多少时间？     输入描述 第一行输出三个正整数，代表城市数量、道路数量以及顺丰小哥使用魔法需要耗费的时间。接下来的行，每行输出三个正整数，代表号城市和号城市有一条双向道路连接，且顺丰小哥使用该道路需要花费的时间为。 输出描述 如果顺丰小哥无法完成运货，则输出-1。否则输出最小需要花费的时间。

# !传送门=>虚拟结点


import heapq
import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m, x = map(int, input().split())
    dist = [[INF, INF] for _ in range(n + 1)]
    magicPoint = n
    adjList = [[] for _ in range(n + 1)]
    for _ in range(m):
        u, v, w = map(int, input().split())
        adjList[u - 1].append((v - 1, 2 * w))
        adjList[v - 1].append((u - 1, 2 * w))
    for i in range(n):
        if i == 0 or i == n - 1:
            continue
        adjList[i].append((magicPoint, x))
        adjList[magicPoint].append((i, x))

    dist[0][0] = 0
    pq = [(0, 0, 0)]  # (dist, node, used)
    while pq:
        curDist, curNode, used = heapq.heappop(pq)
        if dist[curNode][used] < curDist:
            continue
        if curNode == n - 1:
            print(curDist // 2)
            exit(0)
        for nxtNode, nxtDist in adjList[curNode]:
            if nxtNode == magicPoint and used:
                continue
            cand1 = curDist + nxtDist
            if cand1 < dist[nxtNode][used]:
                dist[nxtNode][used] = cand1
                heapq.heappush(pq, (cand1, nxtNode, used))
        # 使用魔法
        if not used:
            cand2 = curDist + x
            if cand2 < dist[curNode][1]:
                dist[curNode][1] = cand2
                heapq.heappush(pq, (cand2, curNode, 1))
    print(-1)
