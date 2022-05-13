from typing import List
from heapq import heappop, heappush

# 如果以最佳方式使用给定的梯子和砖块，返回你可以到达的最远建筑物的下标（下标 从 0 开始 ）。
# 总结：优先使用梯子，将diff保存最小堆，梯子不够了用砖头兑换最小的diff
class Solution:
    def furthestBuilding(self, heights: List[int], bricks: int, ladders: int) -> int:
        pq = []
        for i in range(len(heights) - 1):
            diff = heights[i + 1] - heights[i]
            if diff > 0:
                heappush(pq, diff)
            if len(pq) > ladders:
                bricks -= heappop(pq)
            if bricks < 0:
                return i

        return len(heights) - 1


print(Solution().furthestBuilding(heights=[4, 2, 7, 6, 9, 14, 12], bricks=5, ladders=1))
