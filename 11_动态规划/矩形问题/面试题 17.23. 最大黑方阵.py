# 给定一个方阵，其中每个单元(像素)非黑即白。
# 设计一个算法，找出 4 条边皆为黑色像素的最大子方阵。

# 返回一个数组 [r, c, size] ，其中 r, c 分别代表子方阵左上角的行号和列号，
# size 是子方阵的边长。若有多个满足条件的子方阵，返回 r 最小的，
# 若 r 相同，返回 c 最小的子方阵。若无满足条件的子方阵，返回空数组。
# 0:黑色 1:白色
# ROW,COL<=200
# 11_动态规划/矩形问题/1139. 最大的以 1 为边界的正方形.py

from typing import List, Tuple


class Solution:
    def findSquare(self, matrix: List[List[int]]) -> Tuple[int, int, int]:
        ROW, COL = len(matrix), len(matrix[0])
        up, left = [[0] * COL for _ in range(ROW)], [[0] * COL for _ in range(ROW)]
        for row in range(ROW):
            for col in range(COL):
                if matrix[row][col] == 0:
                    up[row][col] = up[row - 1][col] + 1 if row else 1
                    left[row][col] = left[row][col - 1] + 1 if col else 1

        resRow, resCol, res = 0, 0, 0
        for r in range(ROW):
            for c in range(COL):
                # 从最大的边长开始遍历
                for k in range(min(up[r][c], left[r][c]), res, -1):
                    if up[r][c - k + 1] >= k and left[r - k + 1][c] >= k:
                        resRow, resCol, res = r - k + 1, c - k + 1, k
                        break

        return (resRow, resCol, res) if res else tuple()
