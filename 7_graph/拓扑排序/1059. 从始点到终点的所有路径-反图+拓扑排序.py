# https://leetcode.cn/problems/all-paths-from-source-lead-to-destination/
# 1059. 从始点到终点的所有路径-反图+拓扑排序
# !判断有向图起点s到终点t的所有路径是否都最终结束于终点t.
# !可能存在重边和自环


from collections import deque
from typing import List


class Solution:
    def leadsToDestination(
        self, n: int, edges: List[List[int]], source: int, destination: int
    ) -> bool:
        """判断从起点出发的所有路径是否最终结束于终点."""
        rAdjList = [[] for _ in range(n)]
        indeg = [0] * n
        for u, v in edges:
            rAdjList[v].append(u)
            indeg[u] += 1

        if indeg[destination]:
            return False

        queue = deque([destination])
        while queue:
            cur = queue.popleft()
            if cur == source:
                return True
            for next in rAdjList[cur]:
                indeg[next] -= 1
                if indeg[next] == 0:
                    queue.append(next)
        return False
