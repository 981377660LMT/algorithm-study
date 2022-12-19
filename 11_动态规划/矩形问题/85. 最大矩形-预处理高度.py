from typing import List


def largestRectangleInHistogram(heights: List[int]) -> int:
    """直方图中的最大矩形"""
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


class Solution:
    def maximalRectangle(self, matrix: List[List[str]]) -> int:
        """
        给定一个仅包含 0 和 1 、大小为 rows x cols 的二维二进制矩阵
        找出只包含 1 的最大矩形，并返回其面积。
        1 <= row, cols <= 200

        每一层看作是直方图
        """
        ROW, COL = len(matrix), len(matrix[0])
        res = 0
        heights = [0] * COL
        for row in range(ROW):
            for col in range(COL):
                if matrix[row][col] == "1":
                    heights[col] += 1
                else:
                    heights[col] = 0
            res = max(res, largestRectangleInHistogram(heights))
        return res


assert (
    Solution().maximalRectangle(
        matrix=[
            ["1", "0", "1", "0", "0"],
            ["1", "0", "1", "1", "1"],
            ["1", "1", "1", "1", "1"],
            ["1", "0", "0", "1", "0"],
        ]
    )
    == 6
)
