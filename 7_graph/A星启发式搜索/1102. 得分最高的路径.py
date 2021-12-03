from typing import List
from heapq import heappush, heappop

# 路径的得分是`该路径上的 最小 值`。例如，路径 8 →  4 →  5 →  9 的值为 4
# 找出所有路径中得分 最高 的那条路径，返回其 得分。

# 每一步走最高分那一步，就是最高分路径了，并不是所有的和，要不然肯定是所有分数加起来最大
# 本题每一步做决策的过程实际上可以理解为，在能走的范围内，走使得路径分数尽可能大的点，也就是说可以抽象成在一个优先级队列找最大值的过程。


class Solution:
    def maximumMinimumPath(self, grid: List[List[int]]) -> int:
        if not any(grid):
            return 0
        row, col = len(grid), len(grid[0])

        pq = [(-grid[0][0], 0, 0)]
        visited = [[False] * col for _ in range(row)]

        res = grid[0][0]
        while pq:
            num, r, c = heappop(pq)
            if visited[r][c]:
                continue
            visited[r][c] = True

            res = min(res, -num)
            if (r, c) == (row - 1, col - 1):
                break

            for nr, nc in [(r - 1, c), (r + 1, c), (r, c - 1), (r, c + 1)]:
                if 0 <= nr < row and 0 <= nc < col:
                    heappush(pq, (-grid[nr][nc], nr, nc))

        return res


print(Solution().maximumMinimumPath([[5, 4, 5], [1, 2, 6], [7, 4, 6]]))

