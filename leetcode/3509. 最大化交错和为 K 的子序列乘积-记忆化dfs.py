# 3509. 最大化交错和为 K 的子序列乘积
# https://leetcode.cn/problems/maximum-product-of-subsequences-with-an-alternating-sum-equal-to-k/description/
# 给你一个整数数组 nums 和两个整数 k 与 limit，你的任务是找到一个非空的 子序列，满足以下条件：
#
# 它的 交错和(alternating sum) 等于 k。
# 在乘积 不超过 limit 的前提下，最大化 其所有数字的乘积。
# 返回满足条件的子序列的 乘积 。如果不存在这样的子序列，则返回 -1。
# 1 <= nums.length <= 150
# 0 <= nums[i] <= 12
# -105 <= k <= 105
# 1 <= limit <= 5000
#
# !1.超过 limit 的乘积一律视作 limit+1，减少状态个数
# !2.剪枝.

from functools import lru_cache
from typing import List


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxProduct(self, nums: List[int], k: int, limit: int) -> int:
        if sum(nums) < abs(k):
            return -1

        @lru_cache(None)
        def dfs(pos: int, alterSum: int, mul: int, sign: int, empty: bool) -> int:
            if pos == len(nums):
                if alterSum == k and mul <= limit and not empty:
                    return 1
                return -1
            res1 = dfs(pos + 1, alterSum, mul, sign, empty)
            cur = nums[pos]
            tmp = dfs(
                pos + 1,
                alterSum + sign * cur,
                min2(mul * cur, limit + 1),
                -sign,
                False,
            )
            res2 = tmp * cur if tmp != -1 else -1
            return max2(res1, res2)

        res = dfs(0, 0, 1, 1, True)
        dfs.cache_clear()
        return res
