# 3098. 求出所有子序列的能量和
# https://leetcode.cn/problems/find-the-sum-of-subsequence-powers/description/
# 给你一个长度为 n 的整数数组 nums 和一个 正 整数 k 。
# 一个子序列的 能量 定义为子序列中 任意 两个元素的差值绝对值的 最小值 。
# 请你返回 nums 中长度 等于 k 的 所有 子序列的 能量和 。
# 由于答案可能会很大，将答案对 109 + 7 取余 后返回。
#
# 2 <= n == nums.length <= 50
# -1e8 <= nums[i] <= 1e8
# 2 <= k <= n
#
# !子序列：选或不选
# 为了求差，考虑先排序，然后dfs记录上一个数是谁、当前选了几个数、当前选的数的最小差值.


from functools import lru_cache
from typing import List


MOD = int(1e9 + 7)

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def sumOfPowers(self, nums: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(index: int, remain: int, pre: int, minDiff: int) -> int:
            if remain < 0:
                return 0
            if index == n:
                return minDiff if remain == 0 else 0

            # 不选
            res = dfs(index + 1, remain, pre, minDiff)
            # 选
            res += dfs(index + 1, remain - 1, nums[index], min2(minDiff, nums[index] - pre))
            return res % MOD

        n = len(nums)
        nums = sorted(nums)
        res = dfs(0, k, -INF, INF)
        dfs.cache_clear()
        return res
