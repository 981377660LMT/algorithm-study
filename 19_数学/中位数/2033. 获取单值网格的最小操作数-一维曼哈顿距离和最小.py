# 给你一个大小为 m x n 的二维整数网格 grid 和一个整数 x 。
# 每一次操作，你可以对 grid 中的任一元素 加 x 或 减 x 。
# 单值网格 是全部元素都相等的网格。
# 返回使网格化为单值网格所需的 最小 操作数。如果不能，返回 -1 。
from itertools import pairwise
from typing import List


class Solution:
    def minOperations(self, grid: List[List[int]], x: int) -> int:
        ROW, COL = len(grid), len(grid[0])
        nums = [grid[i][j] for i in range(ROW) for j in range(COL)]
        nums.sort()
        for pre, cur in pairwise(nums):
            if (cur - pre) % x != 0:
                return -1
        mid = nums[len(nums) // 2]
        return sum(abs(num - mid) // x for num in nums)
