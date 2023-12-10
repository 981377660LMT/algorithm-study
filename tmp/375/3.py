from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个整数数组 nums 和一个 正整数 k 。

# 请你统计有多少满足 「 nums 中的 最大 元素」至少出现 k 次的子数组，并返回满足这一条件的子数组的数目。

# 子数组是数组中的一个连续元素序列。


class Solution:
    def countSubarrays(self, nums: List[int], k: int) -> int:
        max_ = max(nums)
        left = 0
        maxCount = 0
        res = 0
        preSum = [0]
        for num in nums:
            preSum.append(preSum[-1] + (num == max_))
        for right, num in enumerate(nums):
            if num == max_:
                maxCount += 1
            while (
                left <= right
                and (nums[left] == max_ and maxCount > k)
                or (nums[left] != max_ and maxCount >= k)
            ):
                if nums[left] == max_:
                    maxCount -= 1
                left += 1
            if preSum[right + 1] - preSum[left] >= k:
                res += left + 1
        return res


# nums = [1,3,2,3,3], k = 2
print(Solution().countSubarrays([1, 3, 2, 3, 3], 2))
