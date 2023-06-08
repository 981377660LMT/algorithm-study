# https://leetcode.cn/problems/minimize-maximum-value-in-a-grid/

# !给定一个包含 `不同` 正整数的 m × n 整数矩阵 grid。
# 必须将矩阵中的每一个整数替换为正整数，且满足以下条件:
# 在替换之后，同行或同列中的每两个元素的 相对 顺序应该保持 不变。
# 替换后矩阵中的 最大 数目应尽可能 小。
# ROW*COL<=1e5


# !dp,维护行列最大值


from collections import defaultdict
from typing import List


class Solution:
    def minScore(self, grid: List[List[int]]) -> List[List[int]]:
        ROW, COL = len(grid), len(grid[0])
        mp = defaultdict(tuple)
        for r, row in enumerate(grid):
            for c, v in enumerate(row):
                mp[v] = (r, c)

        dp = defaultdict(int)  # (x,y)处的值
        rowDp = [0] * ROW  # 当前行的最大值
        colDp = [0] * COL  # 当前列的最大值

        keys = sorted(mp)
        for key in keys:
            r, c = mp[key]
            dp[(r, c)] = max(rowDp[r], colDp[c]) + 1
            rowDp[r] = max(rowDp[r], dp[(r, c)])
            colDp[c] = max(colDp[c], dp[(r, c)])

        return [[dp[(r, c)] for c in range(COL)] for r in range(ROW)]
