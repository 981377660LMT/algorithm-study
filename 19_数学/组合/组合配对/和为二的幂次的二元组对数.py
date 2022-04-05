# 和为二的幂次的二元组对数
from collections import defaultdict


class Solution:
    def solve(self, nums):
        res = 0
        counter = defaultdict(int)
        powers = [(1 << i) for i in range(32)]

        for num in nums:
            for p in powers:
                target = p - num
                res += counter[target]
            counter[num] += 1

        return res
