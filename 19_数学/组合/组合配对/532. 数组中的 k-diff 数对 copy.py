from typing import List
from collections import Counter


class Solution:
    def findPairs(self, nums: List[int], k: int) -> int:
        res = 0
        counter = Counter(nums)
        if k == 0:
            for num in counter:
                res += int(counter[num] >= 2)
        else:
            for num in counter:
                res += int(counter[num + k] >= 1)
        return res

