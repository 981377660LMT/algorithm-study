from collections import defaultdict, deque
from typing import DefaultDict, List, Set


# 无向基环树 拓扑排序 找环
# 无向图拓扑排序，剪掉所有树枝


class Solution:
    def distanceToCycle(self, n: int, edges: List[List[int]]) -> List[int]:
        def findCycle(n: int, adjMap: DefaultDict[int, Set[int]], degrees: List[int]) -> Set[int]:
            """无向图找环上的点 拓扑排序，剪掉所有树枝"""
            queue = deque([i for i in range(n) if degrees[i] == 1])
            visited = [False] * n
            while queue:
                cur = queue.popleft()
                visited[cur] = True
                for next in adjMap[cur]:
                    if visited[next]:
                        continue
                    degrees[next] -= 1
                    if degrees[next] == 1:
                        queue.append(next)

            onCycle = [i for i in range(n) if not visited[i]]
            return set(onCycle)

        """无向图中恰有一个环"""
        adjMap = defaultdict(set)
        degrees = [0] * n
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)
            degrees[u] += 1
            degrees[v] += 1

        cycle = findCycle(n, adjMap, degrees)

        # 从基环出发，求所有树枝上的点的深度
        res = [int(1e20)] * n
        for index in cycle:
            res[index] = 0

        queue = deque([(i, 0) for i in cycle])
        while queue:
            cur, dist = queue.popleft()
            for next in adjMap[cur]:
                if res[next] > dist + 1:
                    res[next] = dist + 1
                    queue.append((next, dist + 1))

        return res
