# 100218. 求出所有子序列的能量和
# https://leetcode.cn/problems/find-the-sum-of-subsequence-powers/description/
# 给你一个长度为 n 的整数数组 nums 和一个 正 整数 k 。
# 一个子序列的 能量 定义为子序列中 任意 两个元素的差值绝对值的 最小值 。
# 请你返回 nums 中长度 等于 k 的 所有 子序列的 能量和 。
# 由于答案可能会很大，将答案对 1e9 + 7 取余 后返回。


# !子序列：选或不选

# 2 <= n == nums.length <= 50
# -1e8 <= nums[i] <= 1e8
# 2 <= k <= n

from typing import List


MOD = int(1e9 + 7)


class Solution:
    def sumOfPowers(self, nums: List[int], k: int) -> int:
        ...
