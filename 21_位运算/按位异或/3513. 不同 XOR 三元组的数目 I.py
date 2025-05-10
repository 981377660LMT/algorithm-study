# 3513. 不同 XOR 三元组的数目 I
# https://leetcode.cn/problems/number-of-unique-xor-triplets-i/description/
# 给你一个长度为 n 的整数数组 nums，其中 nums 是范围 [1, n] 内所有数的 排列 。
# XOR 三元组 定义为三个元素的异或值 nums[i] XOR nums[j] XOR nums[k]，其中 i <= j <= k。
# 返回所有可能三元组 (i, j, k) 中 不同 的 XOR 值的数量。
# 排列 是一个集合中所有元素的重新排列。


from typing import List


class Solution:
    def uniqueXorTriplets(self, nums: List[int]) -> int:
        n = len(nums)
        return n if n <= 2 else 1 << n.bit_length()
