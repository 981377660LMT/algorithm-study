# 和可被K整除的本质不同子数组个数
# https://leetcode.cn/problems/count-distinct-subarrays-divisible-by-k-in-sorted-array/description/
# 给你一个按 非降序 排列的整数数组 nums 和一个正整数 k。
# 如果 nums 的某个 子数组 的元素和可以被 k 整除，则称其为 良好 子数组。
# 返回一个整数，表示 nums 中 不同 的 良好 子数组的数量。
# 子数组 是数组中连续且 非空 的一段元素序列。
#
# !对于连续相同元素段，要保证哈希表暂时不包含这一段对应的前缀和，等我们遍历完这一段，再把对应的前缀和加到哈希表中。

from typing import List
from collections import defaultdict


class Solution:
    def numGoodSubarrays(self, nums: List[int], k: int) -> int:
        counter = defaultdict(int)
        counter[0] = 1
        presum = 0
        res = 0
        ptr = 0
        for i, x in enumerate(nums):
            if i and x != nums[i - 1]:
                v = nums[i - 1]
                s = presum
                for _ in range(i - ptr):
                    counter[s % k] += 1
                    s -= v
                ptr = i
            presum += x
            res += counter[presum % k]
        return res
