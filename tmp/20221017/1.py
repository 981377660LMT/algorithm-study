from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 不包含 任何零的整数数组 nums ，找出自身与对应的负数都在数组中存在的最大正整数 k 。
# 返回正整数 k ，如果不存在这样的整数，返回 -1


class Solution:
    def findMaxK(self, nums: List[int]) -> int:
        S = set(nums)
        res = -1
        for x in nums:
            if x > 0 and -x in S:
                res = max(res, x)
        return res
