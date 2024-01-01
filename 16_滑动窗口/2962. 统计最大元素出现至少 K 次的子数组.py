# 2962. 统计最大元素出现至少 K 次的子数组
# https://leetcode.cn/problems/count-subarrays-where-max-element-appears-at-least-k-times/description/

from typing import List


class Solution:
    def countSubarrays(self, nums: List[int], k: int) -> int:
        max_ = max(nums)
        res, left, maxCount = 0, 0, 0
        for right, num in enumerate(nums):
            if num == max_:
                maxCount += 1
            while left <= right and maxCount >= k:
                if nums[left] == max_:
                    maxCount -= 1
                left += 1
            res += left
        return res
