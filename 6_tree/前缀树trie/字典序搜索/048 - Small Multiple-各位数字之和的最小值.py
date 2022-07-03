# k的倍数中 求各位数字之和的最小值
# 2≤K≤1e5

# 答案不超过k各位数之和
# 对于某个位数和 是否存在一个数是k的倍数 即modk为0
# !即求到modk=0的最短路径 按照字典序搜索 边权就是多加的那个数
from collections import defaultdict
from heapq import heapify, heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

k = int(input())

dist = defaultdict(lambda: int(1e20))
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

    for i in range(10):
        next = (cur * 10 + i) % k
        if dist[cur] + i < dist[next]:
            dist[next] = dist[cur] + i
            heappush(pq, (dist[next], next))
