from typing import List
from MonoQueue import MonoQueue


class Solution:
    def solve(self, nums: List[int]) -> int:
        """
        return the length of the longest sublist where 2 * min(sublist) > max(sublist)
        求最小值两倍大于最大值的最长子数组
        """
        n = len(nums)
        monoQueue = MonoQueue()
        res = 0

        right = 0
        for left in range(n):
            while right < n and (not monoQueue or 2 * monoQueue.min > monoQueue.max):
                res = max(res, right - left)
                monoQueue.append(nums[right])
                right += 1

            if 2 * monoQueue.min > monoQueue.max:
                res = max(res, right - left)

            monoQueue.popleft()

        return res


# The sublist [5, 5, 3, 3] is the longest sublist that meet the criteria since 2 * 3 > 5.
print(Solution().solve(nums=[9, 1, 5, 5, 3, 3]))
