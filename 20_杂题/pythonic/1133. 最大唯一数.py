from typing import List
from collections import Counter


class Solution:
    def largestUniqueNumber(self, nums: List[int]) -> int:
        return max([k for k, v in Counter(nums).items() if v == 1], default=-1)


print(Solution().largestUniqueNumber([5, 7, 3, 9, 4, 9, 8, 3, 1]))

