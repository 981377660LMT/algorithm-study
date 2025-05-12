# https://leetcode.cn/problems/find-x-value-of-array-i/description/
# !求子数组乘积模k为0,1,...,k-1的个数
# !k <= 5.
#
# dp[i+1][x]表示以nums[i]结尾的子数组乘积模k为x的个数.

from typing import List


class Solution:
    def resultArray(self, nums: List[int], k: int) -> List[int]:
        res = [0] * k
        dp = [0] * k
        for v in nums:
            ndp = [0] * k
            ndp[v % k] = 1  # 单独
            for m, c in enumerate(dp):  # 之前的
                ndp[m * v % k] += c
            dp = ndp
            for m, c in enumerate(dp):
                res[m] += c
        return res
