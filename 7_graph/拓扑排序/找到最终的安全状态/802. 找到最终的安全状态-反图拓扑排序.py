from typing import List
from collections import deque


class Solution:
    def eventualSafeNodes(self, graph: List[List[int]]) -> List[int]:
        #  反图
        n = len(graph)
        indegrees = [0] * n
        reversedGraph = [[] for _ in range(n)]
        for cur, nexts in enumerate(graph):
            for next in nexts:
                reversedGraph[next].append(cur)
                indegrees[cur] += 1

        queue = deque([i for i, d in enumerate(indegrees) if d == 0])
        while queue:
            cur = queue.popleft()
            for next in reversedGraph[cur]:
                indegrees[next] -= 1
                if indegrees[next] == 0:
                    queue.append(next)

        return [i for i, d in enumerate(indegrees) if d == 0]


print(Solution().eventualSafeNodes(graph=[[1, 2], [2, 3], [5], [0], [5], [], []]))
# 输出：[2,4,5,6]
