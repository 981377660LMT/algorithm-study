from collections import defaultdict
from heapq import heappop, heappush

#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
#
# @param serialP int一维数组 牛牛占领的p个星球的编号
# @param serialQ int一维数组 牛妹占领的q个星球的编号
# @param path int二维数组 m条隧道，每条隧道有三个数分别是ui,vi,wi。ui,vi分别是隧道的两边星球的编号，wi是它们之间的距离
# @param nn int 星球个数n
# @return int
# 这两个星球的最短距离是多少。
#


class Solution:
    def Length(self, serialP, serialQ, path, nn):
        adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
        for u, v, w in path:
            adjMap[u][v] = min(adjMap[u][v], w)
            adjMap[v][u] = min(adjMap[v][u], w)

        pq = [(0, num) for num in serialP]
        dist = defaultdict(lambda: int(1e20))
        target = set(serialQ)
        while pq:
            curDist, cur = heappop(pq)
            if cur in target:
                return curDist
            for next, weight in adjMap[cur].items():
                if dist[(cur, next)] > curDist + weight:
                    dist[(cur, next)] = curDist + weight
                    heappush(pq, (curDist + weight, next))
        return -1
