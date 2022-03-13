from heapq import heappop, heappush
from typing import List, Tuple
from collections import defaultdict, deque

INF = int(1e20)


class Solution:
    def minimumWeight(self, n: int, edges: List[List[int]], src1: int, src2: int, dest: int) -> int:
        """请你从图中选出一个 边权和最小 的子图，使得从 src1 和 src2 出发，在这个子图中，都 可以 到达 dest"""
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for cur, v, w in edges:
            # 反图
            adjMap[v][cur] = min(adjMap[v][cur], w)

        # 终点开始找
        def dijk(start, end):
            dist = [INF] * n
            dist[start] = 0
            q = [(0, start, [start])]
            while q:
                curDist, cur, path = heappop(q)
                if cur == end:
                    return curDist, path
                for next in adjMap[cur]:
                    if dist[next] > dist[cur] + adjMap[cur][next]:
                        dist[next] = dist[cur] + adjMap[cur][next]
                        heappush(q, (dist[next], next, path + [next]))
            return INF, []

        res = INF

        oneTo2, _ = dijk(src1, src2)
        twoTo1, _ = dijk(src2, src1)
        tTo1, path1 = dijk(dest, src1)
        tTo2, path2 = dijk(dest, src2)

        if tTo1 == INF or tTo2 == INF:
            return -1
        # 减去重合路径  注意并不一定在这条路上
        lastEnd = dest
        for i, (a, b) in enumerate(zip(path1, path2)):
            if a == b:
                lastEnd = a
            else:
                break
        toPath, _ = dijk(dest, lastEnd)
        res = min(res, tTo1 + tTo2) - (toPath)

        if oneTo2 < INF:
            res = min(res, tTo1 + oneTo2)
        if twoTo1 < INF:
            res = min(res, tTo2 + twoTo1)
        return res if res < INF else -1


print(
    Solution().minimumWeight(
        5, [[0, 2, 1], [0, 3, 1], [2, 4, 1], [3, 4, 1], [1, 2, 1], [1, 3, 10]], 0, 1, 4
    )
)

# 5
# [[0,2,1],[0,3,1],[2,4,1],[3,4,1],[1,2,1],[1,3,10]]
# 0
# 1
# 4

# 预期3 输出4
73
