# https://leetcode.cn/problems/longest-square-streak-in-an-array/description/
# 2501. 数组中最长的方波
# 给你一个整数数组 nums 。如果 nums 的子序列满足下述条件，则认为该子序列是一个 方波 ：
# 1.子序列的长度至少为 2 ，并且
# 2.将子序列从小到大排序 之后 ，除第一个元素外，每个元素都是前一个元素的 平方 。
# 返回 nums 中 最长方波 的长度，如果不存在 方波 则返回 -1 。
# 子序列 也是一个数组，可以由另一个数组删除一些或不删除元素且不改变剩余元素的顺序得到。


from typing import List
from collections import defaultdict


class Solution:
    def longestSquareStreak(self, nums: List[int]) -> int:
        dp = defaultdict(int)
        for num in sorted(nums, reverse=True):
            dp[num] = max(dp[num], dp[num**2] + 1)
        max_ = max(dp.values())
        return max_ if max_ > 1 else -1


print(Solution().longestSquareStreak(nums=[2, 3, 5, 6, 7]))
