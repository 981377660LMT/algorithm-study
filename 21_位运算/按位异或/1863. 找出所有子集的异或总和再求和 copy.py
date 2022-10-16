import operator
from functools import reduce
from typing import List


class Solution:
    def subsetXORSum0(self, nums: List[int]) -> int:
        """O(n) 统计每位的贡献

        子集中第i位对结果贡献位count*(1<<i)
        其中count为子集中第i位元素为1的元素个数为奇数的子集数等于偶数个子集数等于 2**(n-1)个
        """
        base = reduce(operator.or_, nums)
        count = 1 << (len(nums) - 1)
        return base * count

    def subsetXORSum1(self, nums: List[int]) -> int:
        res = 0
        dp = [0] * (1 << len(nums))
        for i, num in enumerate(nums):
            for pre in range(1 << i):
                cur = dp[pre] ^ num
                dp[pre | (1 << i)] = cur
                res += cur
        return res

    def subsetXORSum2(self, nums: List[int]) -> int:
        res = 0
        dp = [0]
        for num in nums:
            ndp = []
            for pre in dp:
                cur = pre ^ num
                res += cur
                ndp.append(cur)
            dp += ndp
        return res
