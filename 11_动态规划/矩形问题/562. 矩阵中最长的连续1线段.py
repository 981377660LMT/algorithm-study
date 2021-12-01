# 给定一个01矩阵 M，找到矩阵中最长的连续1线段。这条线段可以是水平的、垂直的、对角线的或者反对角线的。
# 输入:
# [[0,1,1,0],
#  [0,1,1,0],
#  [0,0,0,1]]
# 输出: 3


# 4个方向，每个方向一个专用二维数组
from typing import List


class Solution:
    def longestLine(self, M: List[List[int]]) -> int:
        row = len(M)
        if row == 0:
            return 0
        col = len(M[0])

        horizon = [[0 for _ in range(col)] for _ in range(row)]
        vertic = [[0 for _ in range(col)] for _ in range(row)]
        diag = [[0 for _ in range(col)] for _ in range(row)]
        anti_diag = [[0 for _ in range(col)] for _ in range(row)]

        res = 0
        for r in range(row):
            for c in range(col):
                if M[r][c] == 0:
                    horizon[r][c] = 0
                    vertic[r][c] = 0
                    diag[r][c] = 0
                    anti_diag[r][c] = 0
                else:
                    horizon[r][c] = horizon[r][c - 1] + 1 if c > 0 else 1
                    vertic[r][c] = vertic[r - 1][c] + 1 if r > 0 else 1
                    diag[r][c] = diag[r - 1][c - 1] + 1 if (r > 0 and c > 0) else 1
                    anti_diag[r][c] = anti_diag[r - 1][c + 1] + 1 if (r > 0 and c + 1 < col) else 1

                res = max(res, horizon[r][c], vertic[r][c], diag[r][c], anti_diag[r][c])

        return res


print(Solution().longestLine([[0, 1, 1, 0], [0, 1, 1, 0], [0, 0, 0, 1]]))
