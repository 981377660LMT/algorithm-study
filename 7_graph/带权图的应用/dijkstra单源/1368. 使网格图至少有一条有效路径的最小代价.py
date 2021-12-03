from typing import List
from heapq import heappop, heappush

# 1234对应右左下上
# 一条 有效路径 为从格子 (0,0) 出发，每一步都顺着数字对应方向走，最终在最右下角的格子 (m - 1, n - 1) 结束的路径
# 你可以花费 cost = 1 的代价修改一个格子中的数字，但每个格子中的数字 只能修改一次 。
# 请你返回让网格图至少有一条有效路径的最小代价。

# 看成最短路问题，需要改方向的边为1，不需要的边为0
# new_cost = cost if grid[i][j] == idx + 1 else cost + 1 (移动的方向与当前箭头方向一致)


class Solution:
    def minCost(self, grid: List[List[int]]) -> int:
        m, n = len(grid), len(grid[0])

        # 对应右左下上
        dirs = ((0, 1), (0, -1), (1, 0), (-1, 0))
        dist = [[0x7FFFFFFF] * n for _ in range(m)]
        pq = []
        heappush(pq, (0, 0, 0))

        while pq:
            cost, i, j = heappop(pq)
            if (i, j) == (m - 1, n - 1):
                return cost

            for idx in range(4):
                next_i, next_j = i + dirs[idx][0], j + dirs[idx][1]
                if next_i >= 0 and next_i < m and next_j >= 0 and next_j < n:
                    # 移动的方向与当前箭头方向一致
                    new_cost = cost if grid[i][j] == (idx + 1) else cost + 1
                    if new_cost < dist[next_i][next_j]:
                        dist[next_i][next_j] = new_cost
                        heappush(pq, (new_cost, next_i, next_j))
        return -1


print(Solution().minCost(grid=[[1, 1, 1, 1], [2, 2, 2, 2], [1, 1, 1, 1], [2, 2, 2, 2]]))
# 输出：3
# 解释：你将从点 (0, 0) 出发。
# 到达 (3, 3) 的路径为： (0, 0) --> (0, 1) --> (0, 2) --> (0, 3) 花费代价 cost = 1 使方向向下 --> (1, 3) --> (1, 2) --> (1, 1) --> (1, 0) 花费代价 cost = 1 使方向向下 --> (2, 0) --> (2, 1) --> (2, 2) --> (2, 3) 花费代价 cost = 1 使方向向下 --> (3, 3)
# 总花费为 cost = 3.

