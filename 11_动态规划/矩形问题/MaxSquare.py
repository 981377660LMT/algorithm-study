# 最大正方形(全1的最大正方形)


from typing import Any, List, Tuple


def maxSquare(grid: List[List[Any]]) -> Tuple[int, Tuple[int, int, int, int]]:
    """二维矩形区域中最大的正方形.

    Args:
        grid (List[List[int]]): 二维矩形区域."1"或者1表示有效区域,"0"或者0表示无效区域.

    Returns:
        Tuple[int, Tuple[int, int, int, int]]: 最大正方形的边长和`[r1,r2) x [c1,c2)`区域.
    """
    ROW = len(grid)
    COL = len(grid[0]) if ROW else 0
    res1 = 0
    res2 = (0, 0, 0, 0)
    dp = [0] * COL
    for r, row in enumerate(grid):
        ndp = [0] * COL
        for c, v in enumerate(row):
            if int(v) == 1:
                if c == 0:
                    ndp[c] = 1
                else:
                    min_ = dp[c - 1]
                    if dp[c] < min_:
                        min_ = dp[c]
                    if ndp[c - 1] < min_:
                        min_ = ndp[c - 1]
                    ndp[c] = min_ + 1
                if ndp[c] > res1:
                    res1 = ndp[c]
                    res2 = (r + 1 - res1, r + 1, c + 1 - res1, c + 1)
        dp = ndp
    return res1, res2


if __name__ == "__main__":
    # 221. 最大正方形
    # https://leetcode.cn/problems/maximal-square/
    class Solution:
        def maximalSquare(self, matrix: List[List[str]]) -> int:
            len_, _ = maxSquare(matrix)
            return len_ * len_
