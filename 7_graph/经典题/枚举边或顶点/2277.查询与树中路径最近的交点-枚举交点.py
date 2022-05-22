from collections import defaultdict, deque
from typing import List

# n<=1000


class Solution:
    def closestNode(self, n: int, edges: List[List[int]], query: List[List[int]]) -> List[int]:
        """n^2"""

        def bfs(start: int) -> None:
            queue = deque([(start, 0)])
            visited = set([start])
            while queue:
                cur, curDist = queue.popleft()
                dist[start][cur] = curDist
                for next in adjMap[cur]:
                    if next in visited:
                        continue
                    visited.add(next)
                    queue.append((next, curDist + 1))

        # 枚举最优交叉点
        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        dist = [[int(1e20)] * n for _ in range(n)]
        for start in range(n):
            bfs(start)

        return [
            min(range(n), key=lambda i: dist[i][root1] + dist[i][root2] + dist[i][root3])
            for root1, root2, root3 in query
        ]

