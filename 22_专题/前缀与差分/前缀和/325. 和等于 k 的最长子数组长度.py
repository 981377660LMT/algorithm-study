# 给定一个数组 nums 和一个目标值 k，找到和等于 k 的最长连续子数组长度。
# 如果不存在任意一个符合要求的子数组，则返回 0。
# !325. 和等于 k 的最长子数组长度 / 和为k的最长子数组长度


from typing import List


class Solution:
    def maxSubArrayLen(self, nums: List[int], k: int) -> int:
        preSum = {0: -1}
        res, curSum = 0, 0
        for i, num in enumerate(nums):
            curSum += num
            if curSum - k in preSum:
                res = max(res, i - preSum[curSum - k])
            preSum.setdefault(curSum, i)
        return res


assert Solution().maxSubArrayLen(nums=[1, -1, 5, -2, 3], k=3) == 4
