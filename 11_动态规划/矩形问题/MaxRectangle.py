# 直方图中的最大矩形


from typing import Any, Callable, List, Tuple


def maxRectangleHistogram(heights: List[int], f: Callable[[int, int, int], None]) -> None:
    """直方图中的最大矩形.
    heights: 直方图的高度.
    f: (start, end, height) : 以 [start, end) 为底, 高度为 height 的矩形.
    """
    n = len(heights)
    stack = []  # (index,value)
    for right in range(n + 1):
        rightHeight = heights[right] if right < n else 0
        j = right
        while stack:
            left, leftHeight = stack[-1]
            if leftHeight < rightHeight:
                break
            f(left, right, leftHeight)
            stack.pop()
            j = left
        stack.append((j, rightHeight))


def maxRectangle1(grid: List[List[Any]], f: Callable[[int, int, int, int], None]) -> None:
    """矩阵中的最大矩形.
    grid: 二维矩阵."1"或者1表示有效区域,"0"或者0表示无效区域.
    f: (r1, r2, c1, c2) : `[r1,r2) x [c1,c2)`区域.
    """
    ROW = len(grid)
    if ROW == 0:
        return

    COL = len(grid[0])
    heights = [0] * COL
    zero = [0] * (COL + 1)
    for i, row in enumerate(grid):
        cache = grid[i + 1] if i + 1 != ROW else []
        for c, v in enumerate(row):
            heights[c] = (heights[c] + 1) if (int(v) == 1) else 0
            if i + 1 != ROW:
                zero[c + 1] = zero[c] + (int(cache[c]) == 0)
            else:
                zero[c + 1] = zero[c] + 1

        # maxRectangleHistogram
        stack = []
        for right in range(COL + 1):
            rightHeight = heights[right] if right < COL else 0
            pos = right
            while stack:
                left, leftHeight = stack[-1]
                if leftHeight < rightHeight:
                    break
                if leftHeight and zero[right] - zero[left]:
                    f(i - leftHeight + 1, i + 1, left, right)
                stack.pop()
                pos = left
            stack.append((pos, rightHeight))


def maxRectangle2(grid: List[List[Any]]) -> Tuple[int, Tuple[int, int, int, int]]:
    """矩阵中的最大矩形.
    Args:
        grid: 二维矩阵."1"或者1表示有效区域,"0"或者0表示无效区域.
    Returns:
        (maxArea, maxRect): maxArea: 最大矩形的面积; maxRect: `[r1,r2) x [c1,c2)`区域.
    """

    def f(r1: int, r2: int, c1: int, c2: int) -> None:
        nonlocal maxArea, maxRect
        area = (r2 - r1) * (c2 - c1)
        if area > maxArea:
            maxArea = area
            maxRect = (r1, r2, c1, c2)

    maxArea, maxRect = 0, (0, 0, 0, 0)
    maxRectangle1(grid, f)
    return maxArea, maxRect


if __name__ == "__main__":
    # 84. 柱状图中最大的矩形
    # https://leetcode.cn/problems/largest-rectangle-in-histogram/
    class Solution1:
        def largestRectangleArea(self, heights: List[int]) -> int:
            def f(start: int, end: int, height: int) -> None:
                nonlocal res
                cand = height * (end - start)
                res = cand if cand > res else res

            res = 0
            maxRectangleHistogram(heights, f)
            return res

    # 85. 最大矩形
    # https://leetcode.cn/problems/maximal-rectangle/
    class Solution2:
        def maximalRectangle(self, matrix: List[List[str]]) -> int:
            return maxRectangle2(matrix)[0]
