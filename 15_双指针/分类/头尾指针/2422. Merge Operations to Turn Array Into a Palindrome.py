# 合并相邻项以形成回文的最少操作次数
from collections import deque
from typing import List


class Solution:
    def minimumOperations(self, nums: List[int]) -> int:
        """deque写法"""
        queue, res = deque(nums), 0
        while len(queue) > 1:
            if queue[0] == queue[-1]:
                queue.popleft()
                queue.pop()
            elif queue[0] < queue[-1]:
                cur = queue.popleft()
                queue[0] += cur
                res += 1
            else:
                cur = queue.pop()
                queue[-1] += cur
                res += 1
        return res

    def minimumOperations2(self, nums: List[int]) -> int:
        """双指针写法"""
        n = len(nums)
        left, right = 0, n - 1
        res = 0
        while left < right:
            if nums[left] == nums[right]:
                left += 1
                right -= 1
            elif nums[left] < nums[right]:
                nums[left + 1] += nums[left]
                left += 1
                res += 1
            else:
                nums[right - 1] += nums[right]
                right -= 1
                res += 1
        return res


print(Solution().minimumOperations(nums=[4, 3, 2, 1, 2, 3, 1]))
print(Solution().minimumOperations2(nums=[4, 3, 2, 1, 2, 3, 1]))
