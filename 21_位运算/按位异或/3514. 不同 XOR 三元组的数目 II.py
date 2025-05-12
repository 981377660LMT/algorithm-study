# 3514. 不同 XOR 三元组的数目 II
# https://leetcode.cn/problems/number-of-unique-xor-triplets-ii/description/
# 给你一个整数数组 nums 。
# XOR 三元组 定义为三个元素的异或值 nums[i] XOR nums[j] XOR nums[k]，其中 i <= j <= k。
# 返回所有可能三元组 (i, j, k) 中 不同 的 XOR 值的数量。
#
# O(n(n+U))
#
# !还可以 UlogU 的fwt.


from typing import List
from itertools import combinations


class Solution:
    def uniqueXorTriplets(self, nums: List[int]) -> int:
        s = set(nums)
        xor2 = {x ^ y for x, y in combinations(s, 2)} | {0}
        xor3 = {xy ^ z for xy in xor2 for z in s}
        return len(xor3)
