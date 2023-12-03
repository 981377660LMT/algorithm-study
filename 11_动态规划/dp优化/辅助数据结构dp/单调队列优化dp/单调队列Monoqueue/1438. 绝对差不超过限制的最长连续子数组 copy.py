"""
返回最长连续子数组的长度，
该子数组中的任意两个元素之间的绝对差必须小于或者等于 limit

两个元素之间的绝对差即为距离 即滑窗内的最大值减去最小值要小于等于limit
MonoQueue 维护滑窗内的最值
"""

from typing import List
from MonoQueue import MonoQueue
from sortedcontainers import SortedList


class Solution:
    def longestSubarray(self, nums: List[int], limit: int) -> int:
        queue = MonoQueue()
        res, left, n = 0, 0, len(nums)
        for right in range(n):
            queue.append(nums[right])
            while left <= right and queue.max - queue.min > limit:
                queue.popleft()
                left += 1
            res = max(res, right - left + 1)
        return res

    def longestSubarray2(self, nums: List[int], limit: int) -> int:
        sl = SortedList()
        res, left, n = 0, 0, len(nums)
        for right in range(n):
            sl.add(nums[right])
            while left <= right and sl[-1] - sl[0] > limit:
                sl.remove(nums[left])
                left += 1
            res = max(res, right - left + 1)
        return res


print(Solution().longestSubarray([8, 2, 4, 7], 4))
print(Solution().longestSubarray2([8, 2, 4, 7], 4))
