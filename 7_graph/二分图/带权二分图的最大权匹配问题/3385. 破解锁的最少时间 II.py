# 3385. 破解锁的最少时间 II
# https://leetcode.cn/problems/minimum-time-to-break-locks-ii/description/

from typing import List
from scipy.optimize import linear_sum_assignment


class Solution:
    def findMinimumTime(self, strength: List[int]) -> int:
        n = len(strength)
        # 代价矩阵
        # 行对应锁 i，列对应打开次序 j (j从0开始表示第j+1次打开)
        # cost[i][j] = ceil(strength[i] / (j+1))
        cost = []
        for i in range(n):
            row = []
            for j in range(n):
                a = strength[i]
                b = j + 1
                row.append((a + b - 1) // b)
            cost.append(row)
        row_ind, col_ind = linear_sum_assignment(cost)
        res = sum(cost[i][j] for i, j in zip(row_ind, col_ind))
        return res
