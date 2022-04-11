# Assuming you can group any point a and b if the Euclidean distance between them is ≤ k,
# return the total number of disjoint groups.

# 按照距离给点分组 所有<=k的分到一组 求组数
from collections import deque

# n ≤ 1,000
class Solution:
    def solve(self, points, k):
        def distance(x1, y1, x2, y2):
            return ((x1 - x2) ** 2 + (y1 - y2) ** 2) ** 0.5

        def bfs(x1, y1):
            if (x1, y1) in visited:
                return
            visited.add((x1, y1))
            queue = deque([(x1, y1)])
            while queue:
                x1, y1 = queue.popleft()
                for x2, y2 in points:
                    if not (x2, y2) in visited:
                        if distance(x1, y1, x2, y2) <= k:
                            queue.append((x2, y2))
                            visited.add((x2, y2))

        visited = set()

        count = 0
        # O(n^2)
        for (x, y) in points:
            if not (x, y) in visited:
                bfs(x, y)
                count += 1

        return count
