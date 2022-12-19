# 84. 柱状图(直方图)中最大的矩形
from typing import List


class Solution:
    def largestRectangleArea(self, heights: List[int]) -> int:
        """每个柱子作为非严格最小值的左右边界"""
        n = len(heights)
        leftMost, rightMost = [0] * n, [n - 1] * n

        stack = []
        for i in range(n):
            while stack and heights[stack[-1]] > heights[i]:
                rightMost[stack.pop()] = i - 1
            stack.append(i)

        stack = []
        for i in range(n - 1, -1, -1):
            while stack and heights[stack[-1]] > heights[i]:
                leftMost[stack.pop()] = i + 1
            stack.append(i)

        return max((rightMost[i] - leftMost[i] + 1) * heights[i] for i in range(n))


assert Solution().largestRectangleArea([2, 1, 5, 6, 2, 3]) == 10
