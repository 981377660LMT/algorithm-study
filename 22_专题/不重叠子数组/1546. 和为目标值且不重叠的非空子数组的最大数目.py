from typing import List

# 请你返回 非空不重叠 子数组的最大数目，且每个子数组中数字和都为 target 。
class Solution:
    def maxNonOverlapping(self, nums: List[int], target: int) -> int:
        res, running_sum = 0, 0
        lookup = {0: -1}
        for i, num in enumerate(nums):
            running_sum += num
            if running_sum - target in lookup:
                res += 1
                # 关键
                lookup.clear()
            lookup[running_sum] = i
        return res


print(Solution().maxNonOverlapping([-1, 3, 5, 1, 4, 2, -9], 6))
