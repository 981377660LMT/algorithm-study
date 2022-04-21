from typing import List
from heapq import heappop, heappush

# 1234对应右左下上
# 一条 有效路径 为从格子 (0,0) 出发，每一步都顺着数字对应方向走，最终在最右下角的格子 (m - 1, n - 1) 结束的路径
# 你可以花费 cost = 1 的代价修改一个格子中的数字，但每个格子中的数字 只能修改一次 。
# 请你返回让网格图至少有一条有效路径的最小代价。

# 看成最短路问题，需要改方向的边为1，不需要的边为0
# 0/1最短路其实可以不用pq 用deque就行
# new_cost = cost if grid[i][j] == idx + 1 else cost + 1 (移动的方向与当前箭头方向一致)


# 对应右左下上
DIRS = ((0, 1), (0, -1), (1, 0), (-1, 0))


class Solution:
    def minCost(self, grid: List[List[int]]) -> int:
        row, col = len(grid), len(grid[0])

        dist = [[int(1e20)] * col for _ in range(row)]
        dist[0][0] = 0
        pq = [(0, 0, 0)]

        while pq:
            cost, curR, curC = heappop(pq)
            if dist[curR][curC] < cost:
                continue
            if (curR, curC) == (row - 1, col - 1):
                return cost

            for index in range(4):
                nextR, nextC = curR + DIRS[index][0], curC + DIRS[index][1]
                if 0 <= nextR < row and 0 <= nextC < col:
                    # 移动的方向与当前箭头方向一致
                    nextCost = cost if grid[curR][curC] == (index + 1) else cost + 1
                    if nextCost < dist[nextR][nextC]:
                        dist[nextR][nextC] = nextCost
                        heappush(pq, (nextCost, nextR, nextC))
        return -1


print(Solution().minCost(grid=[[1, 1, 1, 1], [2, 2, 2, 2], [1, 1, 1, 1], [2, 2, 2, 2]]))
# 输出：3
# 解释：你将从点 (0, 0) 出发。
# 到达 (3, 3) 的路径为： (0, 0) --> (0, 1) --> (0, 2) --> (0, 3) 花费代价 cost = 1 使方向向下 --> (1, 3) --> (1, 2) --> (1, 1) --> (1, 0) 花费代价 cost = 1 使方向向下 --> (2, 0) --> (2, 1) --> (2, 2) --> (2, 3) 花费代价 cost = 1 使方向向下 --> (3, 3)
# 总花费为 cost = 3.

