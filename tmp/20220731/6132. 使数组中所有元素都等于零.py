# 给你一个非负整数数组 nums 。在一步操作中，你必须：

# 选出一个正整数 x ，x 需要小于或等于 nums 中 最小 的 非零 元素。
# nums 中的每个正整数都减去 x。
# 返回使 nums 中所有元素都等于 0 需要的 最少 操作数。
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minimumOperations(self, nums: List[int]) -> int:
        res = 0
        while any(num > 0 for num in nums):
            res += 1
            cand = int(1e20)
            for num in nums:
                if num > 0:
                    cand = min(cand, num)
            nums = [num - cand for num in nums]
        return res
