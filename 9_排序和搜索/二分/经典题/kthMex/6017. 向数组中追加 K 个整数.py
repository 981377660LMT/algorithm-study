from typing import List
from bisect import bisect_left

# ps:这道题挺坑的，没看数据范围WA了一次，nums没去重WA了一次😂
# 记得要去重
# 1 <= nums.length <= 105
# 1 <= nums[i], k <= 109


class Solution:
    def minimalKSum(self, nums: List[int], k: int) -> int:
        def findMex(nums: List[int], k: int) -> int:
            """二分搜索缺失的第k个正整数,lc1539. 第 k 个缺失的正整数"""
            # MEX:Min Excluded
            nums = sorted(set(nums))
            left, right = 0, len(nums) - 1
            while left <= right:
                mid = (left + right) >> 1
                diff = nums[mid] - (mid + 1)
                if diff >= k:
                    right = mid - 1
                else:
                    left = mid + 1
            return left + k

        nums = sorted(set(nums))
        mex = findMex(nums, k)
        index = bisect_left(nums, mex)
        allsum = (mex + 1) * (mex) // 2
        return allsum - sum(nums[:index])
