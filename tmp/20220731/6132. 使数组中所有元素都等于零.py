# 给你一个非负整数数组 nums 。在一步操作中，你必须：

# 选出一个正整数 x ，x 需要小于或等于 nums 中 最小 的 非零 元素。
# !nums 中的每个正整数都减去 x。
# 返回使数组中所有元素都等于零需要的 最少 操作数。
from typing import List


class Solution:
    def minimumOperations(self, nums: List[int]) -> int:
        """相当于每次可以把最小的所有元素变为0 即求元素个数"""
        return len(set(nums) - {0})

    def minimumOperations2(self, nums: List[int]) -> int:
        """模拟"""
        res = 0
        while any(num > 0 for num in nums):
            res += 1
            min_ = min(num for num in nums if num > 0)
            nums = [num - min_ for num in nums]
        return res
