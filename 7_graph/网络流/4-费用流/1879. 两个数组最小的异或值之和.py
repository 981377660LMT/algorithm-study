# 请你将 nums2 中的元素重新排列，使得 异或值之和 最小 。

from typing import List
from MinCostMaxFlow import MinCostMaxFlowDinic


class Solution:
    def minimumXORSum(self, nums1: List[int], nums2: List[int]) -> int:
        n = len(nums1)
        START, END = 2 * n + 1, 2 * n + 2
        mcmf = MinCostMaxFlowDinic(END + 1, START, END)
        for i in range(n):
            mcmf.addEdge(START, i, 1, 0)
            mcmf.addEdge(i + n, END, 1, 0)
            for j in range(n):
                mcmf.addEdge(i, j + n, 1, nums1[i] ^ nums2[j])
        return mcmf.work()[1]


assert Solution().minimumXORSum(nums1=[1, 2], nums2=[2, 3]) == 2
