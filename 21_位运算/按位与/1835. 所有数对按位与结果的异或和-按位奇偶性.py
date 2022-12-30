# 按位讨论奇偶性


from typing import List


class Solution:
    def getXORSum(self, arr1: List[int], arr2: List[int]) -> int:
        res = 0
        for bit in range(32, -1, -1):
            count1 = sum((a >> bit) & 1 for a in arr1)
            count2 = sum((a >> bit) & 1 for a in arr2)
            res += (((count1 & 1) & (count2 & 1)) & 1) << bit
        return res


assert Solution().getXORSum(arr1=[1, 2, 3], arr2=[6, 5]) == 0
