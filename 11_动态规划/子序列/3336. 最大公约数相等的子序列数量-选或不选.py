# 3336. 最大公约数相等的子序列数量-选或不选
# https://leetcode.cn/problems/find-the-number-of-subsequences-with-equal-gcd/description/
# 给你一个整数数组 nums。
# 请你统计所有满足一下条件的 非空 子序列对 (seq1, seq2) 的数量：
# 子序列 seq1 和 seq2 不相交，意味着 nums 中 不存在 同时出现在两个序列中的下标。
# seq1 元素的 GCD 等于 seq2 元素的 GCD。
# 返回满足条件的子序列对的总数。
# 由于答案可能非常大，请返回其对 109 + 7 取余 的结果。
# gcd(a, b) 表示 a 和 b 的 最大公约数。
# 子序列 是指可以从另一个数组中删除某些或不删除元素得到的数组，并且删除操作不改变其余元素的顺序。
#
#
# 莫比乌斯反演/倍数容斥：https://leetcode.cn/problems/find-the-number-of-subsequences-with-equal-gcd/solutions/2967084/duo-wei-dppythonjavacgo-by-endlesscheng-5pk3/

from functools import lru_cache
from math import gcd
from typing import List


MOD = int(1e9 + 7)


class Solution:
    def subsequencePairCount(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, gcd1: int, gcd2: int) -> int:
            if index == len(nums):
                return gcd1 == gcd2 != 0
            res1 = dfs(index + 1, gcd1, gcd2)
            res2 = dfs(index + 1, gcd(gcd1, nums[index]), gcd2)
            res3 = dfs(index + 1, gcd1, gcd(gcd2, nums[index]))
            return (res1 + res2 + res3) % MOD

        res = dfs(0, 0, 0)
        dfs.cache_clear()
        return res % MOD
