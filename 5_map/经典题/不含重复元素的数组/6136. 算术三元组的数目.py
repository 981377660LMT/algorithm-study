from itertools import combinations
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始、`严格递增` 的整数数组 nums 和一个正整数 diff 。
# 如果满足下述全部条件，则三元组 (i, j, k) 就是一个 算术三元组 ：
# i < j < k ，
# nums[j] - nums[i] == diff 且
# nums[k] - nums[j] == diff


class Solution:
    def arithmeticTriplets(self, nums: List[int], diff: int) -> int:
        """注意到严格递增 每个元素不一样 因此可以枚举检查每个数"""
        ok = set(nums)
        return sum(num + diff in ok and num - diff in ok for num in nums)

    def arithmeticTriplets2(self, nums: List[int], diff: int) -> int:
        res = 0
        for a, b, c in combinations(nums, 3):
            if b - a == diff and c - b == diff:
                res += 1
        return res


# class Solution:
#     def arithmeticTriplets(self, nums: List[int], diff: int) -> int:
#         s = set(nums)
#         ans = 0
#         for x in s:
#             if x + diff in s and x + diff + diff in s:
#                 ans += 1
#         return ans
