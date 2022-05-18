# m,n<=3e4
from bisect import bisect_left


class Solution:
    def findKthNumber(self, m: int, n: int, k: int) -> int:
        """时间复杂度O(min(m,n)logmn)"""
        # 统计表中不大于mid的个数
        if m > n:  # 优化
            m, n = n, m
        countNGT = lambda mid: sum(min(n, mid // row) for row in range(1, m + 1))
        return bisect_left(range(int(m * n + 10)), k, key=countNGT)

