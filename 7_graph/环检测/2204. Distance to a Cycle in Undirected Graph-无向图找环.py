from collections import defaultdict, deque
from typing import DefaultDict, List, Set

# 1. dfs 找环
class Solution:
    def distanceToCycle(self, n: int, edges: List[List[int]]) -> List[int]:
        def findCycle(n: int, adjMap: DefaultDict[int, Set[int]]) -> List[int]:
            def dfs(cur: int, pre: int) -> bool:
                """环检测，并记录路径"""
                if visited[cur]:
                    return True
                visited[cur] = True
                for next in adjMap[cur]:
                    if next == pre:
                        continue
                    path.append(next)
                    if dfs(next, cur):
                        return True
                    path.pop()
                return False

            res = []
            path = []
            visited = [False] * n

            for i in range(n):
                if visited[i]:
                    continue
                path.append(i)
                if dfs(i, -1):
                    break

            last = path.pop()
            res.append(last)
            while path and path[-1] != last:
                res.append(path.pop())

            return res

        """无向图中恰有一个环"""
        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        cycle = findCycle(n, adjMap)

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
