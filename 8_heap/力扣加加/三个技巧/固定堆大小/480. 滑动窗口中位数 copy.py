# use a SortedList structure, which was implemented using self-balanced tree.
# SortedList enables O(logk) add and O(logk) remove. So the total time complexity is O(nlogk).

from typing import List
from sortedcontainers import SortedList


class Solution:
    def medianSlidingWindow(self, nums: List[int], k: int) -> List[float]:
        list = SortedList()
        res = []
        for i in range(len(nums)):
            list.add(nums[i])
            if len(list) > k:
                list.remove(nums[i - k])
            if len(list) == k:
                median = list[k // 2] if k & 1 else (list[k // 2 - 1] + list[k // 2]) / 2
                res.append(median)
        return res


print(Solution().medianSlidingWindow([1, 3, -1, -3, 5, 3, 6, 7], 3))
