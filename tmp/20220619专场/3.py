from typing import List, Tuple
from collections import defaultdict, Counter

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def findMaxCI(self, nums: List[int]) -> int:
        n = len(nums)
        if n <= 1:
            return n

        dp = 1
        res = 1
        pre = nums[0]
        for num in nums[1:]:
            if num > pre:
                dp += 1
            else:
                dp = 1
            res = max(res, dp)
            pre = num

        return res

