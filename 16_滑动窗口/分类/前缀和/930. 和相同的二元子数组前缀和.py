# 请你统计并返回有多少个和为 goal 的 非空 子数组。
# 子数组 是数组的一段连续部分。
from collections import defaultdict
from typing import List


class Solution:
    def numSubarraysWithSum(self, nums: List[int], goal: int) -> int:
        preSum = defaultdict(int, {0: 1})
        res, curSum = 0, 0
        for num in nums:
            curSum += num
            res += preSum[curSum - goal]
            preSum[curSum] += 1
        return res


print(Solution().numSubarraysWithSum([1, 0, 1, 0, 1], 2))
