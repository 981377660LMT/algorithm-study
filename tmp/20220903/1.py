from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def findSubarrays(self, nums: List[int]) -> bool:
        # 判断是否存在 两个 长度为 2 的子数组且它们的 和 相等
        n = len(nums)
        visited = set()
        for i in range(n - 1):
            sum_ = nums[i] + nums[i + 1]
            if sum_ in visited:
                return True
            visited.add(sum_)
        return False
