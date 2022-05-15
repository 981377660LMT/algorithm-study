kanade 算法两种理解

```Python

from typing import List

class Solution:
    def maxSubArray(self, nums: List[int]) -> int:
        """最大子数组和-取或全不取dp"""
        if len(nums) == 1:
            return nums[0]

        curMax, res = -int(1e20), nums[0]
        for num in nums:
            # 如果curMax为负数，则前面(包括自己)全不取
            curMax = max(curMax + num, num)
            res = max(res, curMax)
        return res

    def maxSubArray2(self, nums: List[int]) -> int:
        """最大子数组和-前缀和"""
        if len(nums) == 1:
            return nums[0]

        curSum, preMin = 0, 0
        res = -int(1e20)
        for num in nums:
            curSum += num
            res = max(res, curSum - preMin)
            preMin = min(preMin, curSum)
        return res

```
