# 562. 矩阵中最长的连续1线段
# 给定一个01矩阵 M，找到矩阵中最长的连续1线段。这条线段可以是水平的、垂直的、对角线的或者反对角线的。
# 输入:
# [[0,1,1,0],
#  [0,1,1,0],
#  [0,0,0,1]]
# 输出: 3

# !ROW*COL<=1e4
# 预处理四个方向连续1的长度

from typing import List


class Solution:
    def longestLine(self, M: List[List[int]]) -> int:
        ROW, COL = len(M), len(M[0])

        horizon = [[0 for _ in range(COL)] for _ in range(ROW)]
        vertic = [[0 for _ in range(COL)] for _ in range(ROW)]
        diag = [[0 for _ in range(COL)] for _ in range(ROW)]
        antiDiag = [[0 for _ in range(COL)] for _ in range(ROW)]

        res = 0
        for r in range(ROW):
            for c in range(COL):
                if M[r][c] == 1:
                    horizon[r][c] = horizon[r][c - 1] + 1 if c > 0 else 1
                    vertic[r][c] = vertic[r - 1][c] + 1 if r > 0 else 1
                    diag[r][c] = diag[r - 1][c - 1] + 1 if (r > 0 and c > 0) else 1
                    antiDiag[r][c] = antiDiag[r - 1][c + 1] + 1 if (r > 0 and c + 1 < COL) else 1

                res = max(res, horizon[r][c], vertic[r][c], diag[r][c], antiDiag[r][c])

        return res


assert (
    Solution().longestLine(
        M=[
            [0, 1, 1, 0],
            [0, 1, 1, 0],
            [0, 0, 0, 1],
        ]
    )
    == 3
)
