from collections import defaultdict, deque
from typing import DefaultDict, List, Set


# 拓扑排序 找环


class Solution:
    def distanceToCycle(self, n: int, edges: List[List[int]]) -> List[int]:
        def findCycle(n: int, adjMap: DefaultDict[int, Set[int]], degrees: List[int]) -> List[int]:
            queue = deque([i for i in range(n) if degrees[i] == 1])
            onCycle = [True] * n
            while queue:
                cur = queue.popleft()
                onCycle[cur] = False
                for next in adjMap[cur]:
                    degrees[next] -= 1
                    if degrees[next] == 1:
                        queue.append(next)

            cycle = [i for i, v in enumerate(onCycle) if v]
            return cycle

        """无向图中恰有一个环"""
        adjMap = defaultdict(set)
        degrees = [0] * n
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)
            degrees[u] += 1
            degrees[v] += 1

        cycle = findCycle(n, adjMap, degrees)

        res = [int(1e20)] * n
        for index in cycle:
            res[index] = 0

        queue = deque([i for i in cycle])
        dist = 0
        while queue:
            length = len(queue)
            for _ in range(length):
                cur = queue.popleft()
                for next in adjMap[cur]:
                    if res[next] > dist + 1:
                        res[next] = dist + 1
                        queue.append(next)
            dist += 1
        return res
