# k的倍数中 求各位数字之和的最小值 (同余最短路)
# 2≤K≤1e5

# 答案不超过k各位数之和
# 反向思考:对于某个位数和 `是否存在一个数是k的倍数` 即modk为0
# !即求到modk=0的最短路径 按照字典序搜索 边权就是新增的digit

from heapq import heappop, heappush
import sys


input = sys.stdin.readline
INF = int(1e18)

if __name__ == "__main__":
    k = int(input())

    adjList = [[] for _ in range(k)]
    for mod in range(k):
        for i in range(10):
            adjList[mod].append(((mod * 10 + i) % k, i))

    dist = [INF] * k
    dist[1] = 1
    pq = [(1, 1)]

    while pq:
        curDist, cur = heappop(pq)
        if curDist > dist[cur]:
            continue
        if cur == 0:  # type: ignore
            print(curDist)
            exit(0)
        for next, weight in adjList[cur]:
            cand = curDist + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
