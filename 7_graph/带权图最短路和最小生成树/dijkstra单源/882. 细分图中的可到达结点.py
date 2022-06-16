from typing import List
from collections import defaultdict
from heapq import heappop, heappush

# edges[i] = [ui, vi, cnti]
# cnti 是将边细分后的新节点总数。注意，cnti == 0 表示边不可细分。
# 现在得到一个新的 细分图 ，请你计算从节点 0 出发，可以到达多少个节点？
# 节点 是否可以到达的判断条件 为：如果节点间距离是 maxMoves 或更少，则视为可以到达；否则，不可到达。

INF = int(1e20)


class Solution:
    def reachableNodes(self, edges: List[List[int]], maxMoves: int, n: int) -> int:
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in edges:
            adjMap[u][v] = w + 1
            adjMap[v][u] = w + 1

        dist = [INF] * n
        dist[0] = 0
        pq = [(0, 0)]
        while pq:
            curDist, cur = heappop(pq)
            if dist[cur] < curDist:
                continue
            for next, weight in adjMap[cur].items():
                cand = weight + dist[cur]
                if cand < dist[next]:
                    dist[next] = cand
                    heappush(pq, (cand, next))

        # 可到达的大端点
        res = sum(dist[i] <= maxMoves for i in range(n))

        # 对每一对大端点，计算可到达的小端点：容斥原理
        for u, v, w in edges:
            w1, w2 = maxMoves - dist[u], maxMoves - dist[v]
            res += max(w1, 0) + max(w2, 0) - max(w1 + w2 - w, 0)

        return res


print(Solution().reachableNodes(edges=[[0, 1, 10], [0, 2, 1], [1, 2, 2]], maxMoves=6, n=3))
# 输出：13
# 解释：边的细分情况如上图所示。
# 可以到达的节点已经用黄色标注出来。
