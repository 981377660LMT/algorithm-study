"""
每个火车站之间有快车、慢车
每个站台处 快车切为慢车花费为0 慢车切为快车花费为cost
坐慢车到下一个城市的花费为A[i] 坐快车到下一个城市的花费为B[i]
求到每个火车站的最小花费
"""

from collections import defaultdict
from heapq import heappop, heappush
from itertools import pairwise
import math
from typing import List


INF = math.inf


class Solution:
    def minimumCosts(self, regular: List[int], express: List[int], expressCost: int) -> List[int]:
        """注意到原图有环，所以通用的方法是 dijkstra 求最短路"""
        n = len(regular)
        OFFSET = 2 * n
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for pre, cur in pairwise(range(n + 1)):
            adjMap[pre][cur] = regular[pre]
            adjMap[pre][pre + OFFSET] = expressCost
            adjMap[pre + OFFSET][pre] = 0
            adjMap[pre + OFFSET][cur + OFFSET] = express[pre]

        dist = defaultdict(lambda: INF, {0: 0})
        pq = [(0, 0)]
        while pq:
            curDist, cur = heappop(pq)
            if curDist > dist[cur]:
                continue
            for next, weight in adjMap[cur].items():
                cand = curDist + weight
                if cand < dist[next]:
                    dist[next] = cand
                    heappush(pq, (cand, next))

        return [min(dist[i], dist[i + OFFSET]) for i in range(1, n + 1)]

    def minimumCosts2(self, A: List[int], B: List[int], C: int) -> List[int]:
        """将每个车站视为一个点(缩点),那么转移就是拓扑序dp了"""

        n, res = len(A), []
        dp1, dp2 = 0, C
        for i in range(n):
            dp1, dp2 = min(dp1 + A[i], dp2 + B[i]), min(dp1 + A[i] + C, dp2 + B[i])
            res.append(min(dp1, dp2))
        return res


print(Solution().minimumCosts(regular=[1, 6, 9, 5], express=[5, 2, 3, 10], expressCost=8))
print(Solution().minimumCosts2(A=[1, 6, 9, 5], B=[5, 2, 3, 10], C=8))
print(Solution().minimumCosts2(A=[11, 5, 13], B=[7, 10, 6], C=3))
