from heapq import heappop, heappush
from itertools import accumulate
from typing import List


INF = int(1e18)

min = lambda x, y: x if x < y else y


class Solution:
    def maxSum(self, grid: List[List[int]], limits: List[int], k: int) -> int:
        for g in grid:
            g.sort(reverse=True)
        row, col = len(grid), len(grid[0])
        pq = []
        for i in range(row):
            if limits[i] > 0:
                heappush(pq, (-grid[i][0], i, 0))
        res = 0
        for _ in range(k):
            if not pq:
                break
            v, i, j = heappop(pq)
            res += -v
            if j + 1 < min(col, limits[i]):
                heappush(pq, (-grid[i][j + 1], i, j + 1))
        return res


# 提交结果：解答错误
# 输入：
# [[3],[9],[1]]
# [1,0,0]
# 1
# 输出：
# 9
# 预期：
# 3

print(Solution().maxSum(grid=[[3], [9], [1]], limits=[1, 0, 0], k=1))
