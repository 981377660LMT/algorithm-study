# 3738. 替换至多一个元素后最长非递减子数组
# https://leetcode.cn/problems/longest-non-decreasing-subarray-after-replacing-at-most-one-element/
from typing import List


class Solution:
    def longestSubarray(self, nums: List[int]) -> int:
        n = len(nums)
        pre = [1] * n
        for i in range(1, n):
            pre[i] = pre[i - 1] + 1 if nums[i] >= nums[i - 1] else 1
        suf = [1] * n
        for i in range(n - 2, -1, -1):
            suf[i] = suf[i + 1] + 1 if nums[i] <= nums[i + 1] else 1

        res = max(max(pre), max(suf)) + 1  # 不替换元素的情况
        res = max(res, suf[0], pre[-1])  # 替换第一个或者最后一个元素的情况
        for i in range(1, n - 1):
            if nums[i - 1] <= nums[i + 1]:
                res = max(res, pre[i - 1] + suf[i + 1] + 1)
        return min(res, n)
