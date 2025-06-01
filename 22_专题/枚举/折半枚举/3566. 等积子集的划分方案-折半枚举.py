# 3566. 等积子集的划分方案-折半枚举
# https://leetcode.cn/problems/partition-array-into-two-equal-product-subsets/description/
# 给你一个整数数组 nums，其中包含的正整数 互不相同 ，另给你一个整数 target。
# 请判断是否可以将 nums 分成两个 非空、互不相交 的 子集 ，并且每个元素必须  恰好 属于 一个 子集，使得这两个子集中元素的乘积都等于 target。
#
# 将nums分为A和B.
# A中选一些数放到第一个集合,乘积为a1,剩下的放到第二个集合,乘积为b1.
# B中选一些数放到第一个集合,乘积为a2,剩下的放到第二个集合,乘积为b2.
# a1*b2=b1*a2，我们可以用一个set维护所有ai/bi的最简分数，有交集则为可行.

from math import gcd, prod
from typing import List, Set, Tuple


class Solution:
    def checkEqualPartitions(self, nums: List[int], target: int) -> bool:
        if prod(nums) != target * target:
            return False

        mid = len(nums) // 2
        set1 = self.calc(nums[:mid], target)
        set2 = self.calc(nums[mid:], target)
        return len(set1 & set2) > 0

    def calc(self, nums: List[int], target: int) -> Set[Tuple[int, int]]:
        res = set()

        def dfs(i: int, a: int, b: int):
            if a > target or b > target:
                return
            if i == len(nums):
                g = gcd(a, b)
                res.add((a // g, b // g))
                return
            dfs(i + 1, a * nums[i], b)
            dfs(i + 1, a, b * nums[i])

        dfs(0, 1, 1)

        return res
