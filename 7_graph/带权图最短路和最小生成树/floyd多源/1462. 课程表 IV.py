# floyd/拓扑排序 n<=100
from collections import deque
from itertools import product
from typing import List

INF = int(1e18)


class Solution:
    def checkIfPrerequisite(
        self, n: int, prerequisites: List[List[int]], queries: List[List[int]]
    ) -> List[bool]:
        """回答课程 uj 是否是课程 vj 的先决条件。"""
        adjList = [[] for _ in range(n)]
        deg = [0] * n
        for cur, pre in prerequisites:
            adjList[pre].append(cur)
            deg[cur] += 1

        parent = [0] * n
        queue = deque([i for i in range(n) if deg[i] == 0])
        while queue:
            cur = queue.popleft()
            for next in adjList[cur]:
                parent[next] |= parent[cur] | (1 << cur)
                deg[next] -= 1
                if deg[next] == 0:
                    queue.append(next)
        return [not not (parent[cur] & (1 << pre)) for cur, pre in queries]

    def checkIfPrerequisite2(
        self, n: int, prerequisites: List[List[int]], queries: List[List[int]]
    ) -> List[bool]:
        """回答课程 uj 是否是课程 vj 的先决条件。"""
        dist = [[False] * n for _ in range(n)]
        for cur, pre in prerequisites:
            dist[pre][cur] = True
        for k, i, j in product(range(n), repeat=3):
            dist[i][j] |= dist[i][k] and dist[k][j]
        return [dist[pre][cur] for cur, pre in queries]
