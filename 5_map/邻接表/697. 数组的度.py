from typing import List
from collections import defaultdict


class Solution:
    def findShortestSubArray(self, nums: List[int]) -> int:
        indexes = defaultdict(list)
        for i, num in enumerate(nums):
            indexes[num].append(i)
        maxCount = max(len(lis) for lis in indexes.values())
        return min(lis[-1] - lis[0] + 1 for lis in indexes.values() if len(lis) == maxCount)


print(Solution().findShortestSubArray(nums=[1, 2, 2, 3, 1, 4, 2]))
# 解释：
# 数组的度是 3 ，因为元素 2 重复出现 3 次。
# 所以 [2,2,3,1,4,2] 是最短子数组，因此返回 6 。
