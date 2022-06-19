from itertools import combinations
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


print(Solution().maximumDetonation(bombs=[[2, 1, 3], [6, 1, 4]]))
