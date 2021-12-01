from typing import List
from itertools import accumulate


class Solution:
    def numOfSubarrays(self, arr: List[int], k: int, threshold: int) -> int:
        pre_sum = [0] + list(accumulate(arr))
        target = k * threshold
        return sum(pre_sum[i] - pre_sum[i - k] >= target for i in range(k, len(arr) + 1))


print(Solution().numOfSubarrays(arr=[2, 2, 2, 2, 5, 5, 5, 8], k=3, threshold=4))
