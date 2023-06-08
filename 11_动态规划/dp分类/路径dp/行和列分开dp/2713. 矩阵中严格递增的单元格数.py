# https://leetcode.cn/problems/maximum-strictly-increasing-cells-in-a-matrix/
# 6456. 矩阵中严格递增的单元格数
# https://atcoder.jp/contests/abc224/tasks/abc224_e
# E - Integers on Grid (爬山)
# 有n个方格上有正整数，其余的方格上的数字为0
# 在方格上走，每次只能走同一行或同一列，
# !且到达的格点上的值必须严格大于(真に大きい)当前值，
# 问从给定的n个点出发 最多走多少步
# !ROW,COL<=2e5 n<=2e5

# !注意到只能前往严格大于当前单元格的值，因此按照单元格的值排序即构成DAG，可以进行DP
# !1.倒着考虑 求最长路
# !2.因为相等不能取,所以相等的点`一起处理`
# !3.行和列分开dp rowDp表示当前行的数最多走了多少步 colDp表示当前列的数最多走了多少步
# !dp[(row, col)] = max(rowDp[row], colDp[col]) + 1
# !rowDp[row] = max(rowDp[row], dp[(row, col)])
# !colDp[col] = max(colDp[col], dp[(row, col)])

from collections import defaultdict
from typing import List


INF = int(4e18)


class Solution:
    def maxIncreasingCells(self, mat: List[List[int]]) -> int:
        ROW, COL = len(mat), len(mat[0])
        mp = defaultdict(list)
        for i, row in enumerate(mat):
            for j, num in enumerate(row):
                mp[num].append((i, j))  # 相同值的元素一起处理

        dp = defaultdict(int)  # (x,y)处的最长路
        rowDp = [0] * ROW  # 当前行的数最多走了多少步
        colDp = [0] * COL  # 当前列的数最多走了多少步

        keys = sorted(mp)
        for key in keys:
            pos = mp[key]
            for x, y in pos:
                dp[(x, y)] = max(rowDp[x], colDp[y]) + 1
            for x, y in pos:
                rowDp[x] = max(rowDp[x], dp[(x, y)])
                colDp[y] = max(colDp[y], dp[(x, y)])
        return max(dp.values(), default=0)
