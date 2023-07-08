"""
如果数组中的某个子数组 恰好 只存在 一 个值为 1 的元素，
则认为该子数组是一个 好子数组 。
请你统计将数组 nums 划分成若干 好子数组 的方法数，并以整数形式返回。
由于数字可能很大，返回其对 1e9 + 7 取余 之后的结果。
"""

from typing import List

MOD = int(1e9 + 7)


class Solution:
    def numberOfGoodSubarraySplits(self, nums: List[int]) -> int:
        ones = [i for i, num in enumerate(nums) if num == 1]
        if not ones:
            return 0
        res = 1
        for i in range(1, len(ones)):
            res *= ones[i] - ones[i - 1]
            res %= MOD
        return res
