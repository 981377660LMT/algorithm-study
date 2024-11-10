from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def hasIncreasingSubarrays(self, nums: List[int], k: int) -> bool:
        n = len(nums)
        for v in range(n - 2 * k + 1):
            arr1 = nums[v : v + k]
            arr2 = nums[v + k : v + 2 * k]
            if all(arr1[i] < arr1[i + 1] for i in range(k - 1)) and all(
                arr2[i] < arr2[i + 1] for i in range(k - 1)
            ):
                return True
        return False
