from itertools import combinations
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的 m x n 二进制矩阵 mat 和一个整数 cols ，表示你需要选出的列数。
# 如果一行中，所有的 1 都被你选中的列所覆盖，那么我们称这一行 被覆盖 了。
# 请你返回在选择 cols 列的情况下，被覆盖 的行数 最大 为多少。


class Solution:
    def maximumRows(self, mat: List[List[int]], cols: int) -> int:
        ROW, COL = len(mat), len(mat[0])
        rowState = [0] * ROW
        for i in range(ROW):
            for j in range(COL):
                if mat[i][j] == 1:
                    rowState[i] |= 1 << j

        res = 0
        for select in combinations(range(COL), cols):
            state = 0
            for s in select:
                state |= 1 << s
            cand = sum(rowState[i] & state == rowState[i] for i in range(ROW))
            res = max(res, cand)
        return res


print(Solution().maximumRows(mat=[[0, 0, 0], [1, 0, 1], [0, 1, 1], [0, 0, 1]], cols=2))
