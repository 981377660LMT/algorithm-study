# 3739. 统计主要元素子数组数目 II
# https://leetcode.cn/problems/count-subarrays-with-majority-element-ii/description/
# 给你一个整数数组 nums 和一个整数 target。
# 返回数组 nums 中满足 target 是 主要元素 的 子数组 的数目。
# 一个子数组的 主要元素 是指该元素在该子数组中出现的次数 严格大于 其长度的 一半 。
# 子数组 是数组中的一段连续且 非空 的元素序列。

from typing import List


class Solution:
    def countMajoritySubarrays(self, nums: List[int], target: int) -> int:
        n = len(nums)
        curSum = n
        counter = [0] * (2 * n + 1)
        counter[curSum] = 1
        res, leftSmaller = 0, 0
        for v in nums:
            if v == target:
                leftSmaller += counter[curSum]
                curSum += 1
            else:
                curSum -= 1
                leftSmaller -= counter[curSum]
            res += leftSmaller
            counter[curSum] += 1
        return res
