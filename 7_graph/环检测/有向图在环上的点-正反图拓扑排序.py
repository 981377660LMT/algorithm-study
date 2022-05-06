from typing import DefaultDict, List, Set
from collections import deque


# 判断哪些点在有向图的环上
# https://leetcode-cn.com/problems/find-eventual-safe-states/
class Solution:
    def onCycle(self, graph: List[List[int]]) -> List[int]:
        def topoSort(adjList: List[Set[int]], indegrees: List[int]) -> Set[int]:
            queue = deque([i for i, d in enumerate(indegrees) if d == 0])
            while queue:
                cur = queue.popleft()
                for next in adjList[cur]:
                    indegrees[next] -= 1
                    if indegrees[next] == 0:
                        queue.append(next)

            return set(i for i, d in enumerate(indegrees) if d != 0)

        #  正反图 + 拓扑排序
        n = len(graph)
        adjList = [set() for _ in range(n)]
        rAdjList = [set() for _ in range(n)]
        indegrees = [0] * n
        rIndegrees = [0] * n
        for cur, nexts in enumerate(graph):
            for next in nexts:
                adjList[cur].add(next)
                rAdjList[next].add(cur)
                indegrees[next] += 1
                rIndegrees[cur] += 1

        return list(topoSort(adjList, indegrees) & topoSort(rAdjList, rIndegrees))


print(Solution().onCycle(graph=[[1, 2], [2, 3], [5], [0], [5], [], []]))
# 输出：[2,4,5,6]
