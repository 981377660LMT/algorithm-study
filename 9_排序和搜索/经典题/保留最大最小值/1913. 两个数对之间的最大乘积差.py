from typing import List


# 4 <= nums.length <= 104
# 1 <= nums[i] <= 104
class Solution:
    def maxProductDifference(self, nums: List[int]) -> float:
        L1, L2 = float('-inf'), float('-inf')
        S1, S2 = float('inf'), float('inf')
        for n in nums:
            if n > L1:
                L2 = L1
                L1 = n

            elif n > L2:
                L2 = n

            if n < S1:
                S2 = S1
                S1 = n
            elif n < S2:
                S2 = n

        return L1 * L2 - S1 * S2

