from bisect import bisect_right
from typing import List


class Solution:
    def matrixMedian(self, grid: List[List[int]]) -> int:
        """求每行递增的矩阵的中位数 O(ROW*log(COL)*log(k))"""

        def countNGT(mid: int) -> int:
            res = 0
            for row in grid:
                res += bisect_right(row, mid)
            return res

        ROW, COL = len(grid), len(grid[0])
        target = (ROW * COL) // 2 + 1
        left, right = 0, int(1e6 + 10)
        while left <= right:
            mid = (left + right) // 2
            if countNGT(mid) < target:
                left = mid + 1
            else:
                right = mid - 1
        return left
