from typing import List

# 1 <= m, n <= 300


class Solution:
    def removeOnes(self, grid: List[List[int]]) -> bool:
        xor_ = [num ^ 1 for num in grid[0]]
        return all(row in (grid[0], xor_) for row in grid)


print(Solution().removeOnes(grid=[[0, 1, 0], [1, 0, 1], [0, 1, 0]]))
