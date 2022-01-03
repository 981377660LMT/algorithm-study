from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList
from bisect import bisect_left, bisect_right
from functools import lru_cache
from itertools import accumulate, groupby, combinations
from math import gcd

MOD = int(1e9 + 7)
INF = 0x7FFFFFFF


class Solution:
    def maximumInvitations(self, favorite: List[int]) -> int:
        def dfs(cur: int, count: int) -> None:
            if visited[cur]:
                return
            visited[cur] = True
            self.res = max(self.res, count)
            for next in adjMap[cur]:
                if visited[next]:
                    continue
                dfs(next, count + int(next not in s))

        def getBiEdge():
            nonlocal adjMap
            res = []
            for cur, nei in adjMap.items():
                for next in nei:
                    if cur < next and cur in adjMap[next]:
                        res.append((cur, next))
            return res

        # build graph
        n = len(favorite)
        indegree = [0] * n
        adjMap = defaultdict(set)
        fromMap = defaultdict(set)
        for cur, next in enumerate(favorite):
            adjMap[cur].add(next)
            indegree[next] += 1
            fromMap[next].add(cur)

        biEdge = getBiEdge()
        s = set(p for pair in biEdge for p in pair)
        for u, v in biEdge:
            for f in fromMap[v]:
                if f == u:
                    continue
                adjMap[f].add(u)
                adjMap[u].add(f)
                indegree[u] += 1
                indegree[f] += 1

            for f in fromMap[u]:
                if f == v:
                    continue
                adjMap[f].add(v)
                adjMap[v].add(f)
                indegree[v] += 1
                indegree[f] += 1

        print(adjMap, indegree)

        # topo
        queue = deque([i for i in range(n) if indegree[i] == 0])
        while queue:
            cur = queue.popleft()
            for next in adjMap[cur]:
                indegree[next] -= 1
                if indegree[next] == 0:
                    queue.append(next)

        # maxCycle
        cycle = [p for p, degree in enumerate(indegree) if degree != 0]
        visited = [0] * n
        self.res = 0
        for start in cycle:
            dfs(start, int(start not in s))

        return self.res + len(biEdge) * 2


print(Solution().maximumInvitations(favorite=[2, 2, 1, 2]))
print(Solution().maximumInvitations(favorite=[1, 2, 0]))
print(Solution().maximumInvitations(favorite=[3, 0, 1, 4, 1]))
print(Solution().maximumInvitations(favorite=[1, 0, 0, 2, 1, 4, 7, 8, 9, 6, 7, 10, 8]))
print(Solution().maximumInvitations(favorite=[1, 0, 3, 2, 5, 6, 7, 4, 9, 8, 11, 10, 11, 12, 10]))
