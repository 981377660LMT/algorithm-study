# 高桥(0)的高桥数为0
# 某位研究者与高桥数为n的人在一组 与小于n的研究者不在一组 那么高桥数为n的人的高桥数为n+1
# 求每个人的高桥数 如果不存在 则为-1
# n,m<=1e5
# 所有组的人数统计<=1e5

# !等价:高桥数为源点到各个点的最短距离
# !如果按照描述组内连边建图 会TLE(边数为n^2级别) 需要给每个组加一个虚拟点连接(边数降为1e5级别) 边权为0.5
# 然后bfs即可

# 等效边替代减少边的数量
from collections import defaultdict, deque
from itertools import count
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n, m = map(int, input().split())
virtual = count(int(1e7))
adjMap = defaultdict(set)
for _ in range(m):
    k = int(input())
    group = [int(num) - 1 for num in input().split()]
    center = next(virtual)  # !注意这里的建图 虚拟源点在中心 连接其他点 边权为0.5 最后统一除以2
    for i in group:
        adjMap[i].add(center)
        adjMap[center].add(i)

# bfs
dist = defaultdict(lambda: int(1e9), {0: 0})
queue = deque([(0, 0)])
while queue:
    cur, curDist = queue.popleft()
    for next in adjMap[cur]:
        if dist[next] > curDist + 1:
            dist[next] = curDist + 1
            queue.append((next, curDist + 1))

for i in range(n):
    if dist[i] == int(1e9):
        print(-1)
    else:
        print(dist[i] // 2)

