from collections import defaultdict
from heapq import heappop, heappush
from typing import List

from typing import List, Mapping
from collections import defaultdict
from heapq import heappop, heappush


INF = int(1e20)


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


# 原来是要枚举中间点 想的太少 还去求最短路的公共路径了
# 得先想一下`整体的最短路是啥形状的` 然后倒推
# 这题和开会那题一样吃了没想形状的亏
# 举反例验证
class Solution:
    def minimumWeight(self, n: int, edges: List[List[int]], src1: int, src2: int, dest: int) -> int:
        """
        请你从图中选出一个 边权和最小 的子图，
        使得从 src1 和 src2 出发，在这个子图中，都 可以 到达 dest 。
        如果这样的子图不存在，请返回 -1 。

        结论：必定是经过图中的某个中间点，再汇聚到dest，枚举中间点即可
        dijkstra只能求单源最短路径，就是从一个点出发，到所有点的。
        所以想知道所有点到这个点的最短路径，就要反过来
        """
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        rAdjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in edges:
            adjMap[u][v] = min(adjMap[u][v], w)
            rAdjMap[v][u] = min(rAdjMap[v][u], w)

        dist1 = dijkstra(n, adjMap, src1)
        dist2 = dijkstra(n, adjMap, src2)
        dist3 = dijkstra(n, rAdjMap, dest)

        # 这里可以直接int(1e99) int(1e88)
        res = INF
        for mid in range(n):
            res = min(res, dist1[mid] + dist2[mid] + dist3[mid])
        return res if res < INF else -1
