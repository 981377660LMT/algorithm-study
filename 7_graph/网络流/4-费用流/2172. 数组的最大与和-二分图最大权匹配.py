from typing import List
from MinCostMaxFlow import MinCostMaxFlowDinic


class Solution:
    def maximumANDSum(self, nums: List[int], numSlots: int) -> int:
        """最大费用最大流"""
        OFFSET, START, END = 100, 201, 202
        mcmf = MinCostMaxFlowDinic(203, START, END)

        for i, num in enumerate(nums):
            for j in range(numSlots):
                mcmf.addEdge(i, j + OFFSET, 1, -(num & (j + 1)))

        for i in range(len(nums)):
            mcmf.addEdge(START, i, 1, 0)
        for j in range(numSlots):
            mcmf.addEdge(j + OFFSET, END, 2, 0)  # !每个篮子 至多 有 2 个整数

        return -mcmf.work()[1]
