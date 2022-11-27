# 给你一个整数数组 nums 和一个整数 k ，请你统计并返回 该数组中和为 k 的连续子数组的个数 。
from collections import defaultdict
from typing import List


class Solution:
    def subarraySum(self, nums: List[int], k: int) -> int:
        preSum = defaultdict(int, {0: 1})  # 如果记录索引就是{0: -1}
        res, curSum = 0, 0
        for i, num in enumerate(nums):
            curSum += num
            if curSum - k in preSum:
                res += preSum[curSum - k]
            preSum[curSum] += 1
        return res


assert Solution().subarraySum(nums=[1, 1, 1], k=2) == 2
