# 100367. 切蛋糕的最小总开销 II
# https://leetcode.cn/problems/minimum-cost-for-cutting-cake-ii/
# 有一个 m x n 大小的矩形蛋糕，需要切成 1 x 1 的小块。
# 给你整数 m ，n 和两个数组：
# horizontalCut 的大小为 m - 1 ，其中 horizontalCut[i] 表示沿着水平线 i 切蛋糕的开销。
# verticalCut 的大小为 n - 1 ，其中 verticalCut[j] 表示沿着垂直线 j 切蛋糕的开销。
# 一次操作中，你可以选择任意不是 1 x 1 大小的矩形蛋糕并执行以下操作之一：
# 沿着水平线 i 切开蛋糕，开销为 horizontalCut[i] 。
# 沿着垂直线 j 切开蛋糕，开销为 verticalCut[j] 。
# 每次操作后，这块蛋糕都被切成两个独立的小蛋糕。
# 每次操作的开销都为最开始对应切割线的开销，并且不会改变。
# 请你返回将蛋糕全部切成 1 x 1 的蛋糕块的 最小 总开销。
# 1 <= m, n <= 1e5
# horizontalCut.length == m - 1
# verticalCut.length == n - 1
# 1 <= horizontalCut[i], verticalCut[i] <= 1e3
#
# 原问题转换为：
# 合并两个数组，第i个数的代价为：这个数之前

from typing import List
from mergeTwoArrayWithMinCost import mergeTwoArrayWithMinCost


class Solution:
    def minimumCost(self, m: int, n: int, horizontalCut: List[int], verticalCut: List[int]) -> int:
        resCost, _ = mergeTwoArrayWithMinCost(horizontalCut, verticalCut)
        return resCost
