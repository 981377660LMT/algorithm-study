from itertools import combinations, product
from math import dist
from typing import List
from collections import deque, defaultdict

# 给你数组 bombs ，请你返回在引爆 一个 炸弹的前提下，最多 能引爆的炸弹数目。

# 有向图

# 1 <= bombs.length <= 100


class Solution:
    def maximumDetonation(self, bombs: List[List[int]]) -> int:
        def bfs(cur: int):
            res = 0
            queue = deque([cur])
            visited = set([cur])

            while queue:
                cur = queue.popleft()
                res += 1
                for next in adjMap[cur]:
                    if next not in visited:
                        visited.add(next)
                        queue.append(next)
            return res

        n = len(bombs)
        adjMap = defaultdict(set)
        for i in range(n):
            x1, y1, r1 = bombs[i]
            for j in range(n):
                if i == j:
                    continue
                x2, y2, _ = bombs[j]
                if (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) <= r1 * r1:
                    adjMap[i].add(j)

        return max((bfs(i) for i in range(n)), default=1)

    def maximumDetonation2(self, bombs: List[List[int]]) -> int:
        n = len(bombs)
        adjMap = {i: set([i]) for i in range(n)}
        for i, j in combinations(range(n), 2):
            p1, p2, r1, r2 = bombs[i][:2], bombs[j][:2], bombs[i][2], bombs[j][2]
            d = dist(p1, p2)
            if d <= r1:
                adjMap[i].add(j)
            if d <= r2:
                adjMap[j].add(i)

        for i, j in combinations(range(n), 2):
            if i in adjMap[j]:
                adjMap[j] |= adjMap[i]
            if j in adjMap[i]:
                adjMap[i] |= adjMap[j]
        return max((len(adjMap[i]) for i in range(n)), default=1)


print(Solution().maximumDetonation(bombs=[[2, 1, 3], [6, 1, 4]]))
print(Solution().maximumDetonation2(bombs=[[2, 1, 3], [6, 1, 4]]))
