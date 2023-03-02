# k的倍数中 求各位数字之和的最小值 (同余最短路)
# 2≤K≤1e5

# 答案不超过k各位数之和
# 反向思考:对于某个位数和 `是否存在一个数是k的倍数` 即modk为0
# !即求到modk=0的最短路径 按照字典序搜索 边权就是新增的digit
from collections import defaultdict
from heapq import heapify, heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)
INF = int(1e18)

if __name__ == "__main__":
    k = int(input())

    dist = defaultdict(lambda: INF)
    for start in range(1, 10):
        dist[start] = start
    pq = [(start % k, start) for start in range(1, 10)]
    heapify(pq)

    while pq:
        curDist, cur = heappop(pq)
        if curDist > dist[cur]:
            continue

        if cur == 0:
            print(curDist)
            exit(0)

        for digit in range(10):
            next = (cur * 10 + digit) % k  # 倍数建边
            if dist[cur] + digit < dist[next]:
                dist[next] = dist[cur] + digit
                heappush(pq, (dist[next], next))
