"""
和大于0/和小于0 的最长子数组

看到子数组和 想到前缀和 
看到子数组长度 想到哈希表存索引
"""

from itertools import accumulate
from typing import List


class Solution:
    def longestWPI(self, hours: List[int]) -> int:
        """求子数组和大于零的最大子数组长度 O(n)"""
        hours = [1 if h > 8 else -1 for h in hours]
        preSum = {0: -1}
        res, curSum = 0, 0
        for i, num in enumerate(hours):
            curSum += num
            if curSum > 0:
                res = i + 1
            if curSum - 1 in preSum:
                res = max(res, i - preSum[curSum - 1])
            preSum.setdefault(curSum, i)
        return res

    def longestWPI2(self, hours: List[int]) -> int:
        """求子数组和大于零的最长子数组长度 O(nlogn)"""

        def check(mid: int) -> bool:
            """是否存在长至少为mid的子数组和大于零"""
            for i in range(mid, n + 1):
                if preSum[i] - preMin[i - mid] > 0:
                    return True
            return False

        n = len(hours)
        nums = [1 if h > 8 else -1 for h in hours]
        preSum = list(accumulate(nums, initial=0))
        preMin = list(accumulate(preSum[1:], initial=0, func=min))  # 前缀最小值

        res = 0
        left, right = 0, n
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                res = mid
                left = mid + 1
            else:
                right = mid - 1
        return res


print(Solution().longestWPI2([9, 9, 6, 0, 6, 6, 9]))
print(Solution().longestWPI2([9, 6, 9]))
