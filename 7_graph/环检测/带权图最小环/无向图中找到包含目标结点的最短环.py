from collections import deque
from typing import List

# n, m ≤ 250
# 从目标出发,bfs,回到自己就是最短环


class Solution:
    def solve(self, graph: List[List[int]], target: int) -> int:
        queue = deque([(target, 0)])
        visited = set()
        while queue:
            len_ = len(queue)
            for _ in range(len_):
                cur, cost = queue.popleft()
                if cur in visited:
                    continue
                visited.add(cur)

                for next in graph[cur]:
                    if next == target:
                        return cost + 1
                    queue.append((next, cost + 1))

        return -1


print(Solution().solve(graph=[[1], [2], [0]], target=0))
