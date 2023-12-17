# https://leetcode.cn/problems/maximum-spending-after-buying-items/
# 2931. 购买物品的最大开销
#
# 给定一个每行都是非递增数组的二维矩阵
# 每一天，你可以在一个商店里购买一件物品。具体来说，在第 d 天，你可以：
# !选择商店 i ，购买数组中最右边的物品 j ，开销为 values[i][j] * d 。
# 请你返回购买所有 m * n 件物品需要的 最大开销 。
# 1 <= m == values.length <= 10
# 1 <= n == values[i].length <= 104
# 1 <= values[i][j] <= 106
# values[i] 按照非递增顺序排序。


from typing import List
from heapq import heapify, heappop, heappush


class Solution:
    def maxSpending(self, values: List[List[int]]) -> int:
        ROW, COL = len(values), len(values[0])
        pq = []  # (value, row, col)
        for r, row in enumerate(values):
            pq.append((row[-1], r, COL - 1))
        heapify(pq)
        res = 0
        for day in range(1, ROW * COL + 1):
            min_, r, c = heappop(pq)
            res += min_ * day
            if c > 0:
                heappush(pq, (values[r][c - 1], r, c - 1))
        return res


assert Solution().maxSpending(values=[[8, 5, 2], [6, 4, 1], [9, 7, 3]]) == 285
