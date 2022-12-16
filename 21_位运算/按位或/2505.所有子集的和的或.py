"""
2505. Bitwise OR of All Subsequence Sums
所有子集的和的或
# !n<=1e5 nums[i]<=1e9
"""

from functools import reduce
from itertools import accumulate
from typing import List


class Solution:
    def subsequenceSumOr(self, nums: List[int]) -> int:
        # !结论:答案为所有输入数字、以及输入数字的所有前缀和的或
        preSum = [0] + list(accumulate(nums))
        return reduce(lambda x, y: x | y, nums + preSum)

    def subsequenceSumOr2(self, nums: List[int]) -> int:
        """对每一位进行判断,最后的或能不能取到这一位"""
        bitCounter = [0] * 64
        for num in nums:
            for i in range(32):
                if num & (1 << i):
                    bitCounter[i] += 1
        for i in range(len(bitCounter) - 1):
            bitCounter[i + 1] += bitCounter[i] // 2
        return sum(1 << i for i in range(len(bitCounter)) if bitCounter[i])


assert Solution().subsequenceSumOr(nums=[2, 1, 0, 3]) == 7
# 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7
