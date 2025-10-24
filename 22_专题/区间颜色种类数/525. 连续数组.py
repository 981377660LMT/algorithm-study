# 525. 连续数组
# https://leetcode.cn/problems/contiguous-array/
# 给定一个二进制数组 nums , 找到含有相同数量的 0 和 1 的最长连续子数组，并返回该子数组的长度。

from typing import List


class Solution:
    def findMaxLength(self, nums: List[int]) -> int:
        pos = {0: -1}
        res = 0
        curSum = 0
        for i, v in enumerate(nums):
            curSum += 1 if v == 1 else -1
            if curSum in pos:
                res = max(res, i - pos[curSum])
            else:
                pos[curSum] = i
        return res
