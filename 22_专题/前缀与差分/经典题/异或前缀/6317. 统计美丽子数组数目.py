# 异或前缀/前缀异或
# 给你一个下标从 0 开始的整数数组nums 。每次操作中，你可以：

# 选择两个满足 0 <= i, j < nums.length 的不同下标 i 和 j 。
# 选择一个非负整数 k ，满足 nums[i] 和 nums[j] 在二进制下的第 k 位（下标编号从 0 开始）是 1 。
# 将 nums[i] 和 nums[j] 都减去 2**k 。
# 如果一个子数组内执行上述操作若干次后，该子数组可以变成一个全为 0 的数组，那么我们称它是一个 美丽 的子数组。

# !请你返回数组 nums 中 美丽子数组 的数目。
# 子数组是一个数组中一段连续 非空 的元素序列。


from typing import List
from collections import defaultdict


# !异或为0的子数组个数
class Solution:
    def beautifulSubarrays(self, nums: List[int]) -> int:
        preSum = defaultdict(int, {0: 1})
        res, curXor = 0, 0
        for _, num in enumerate(nums):
            curXor ^= num
            res += preSum[curXor]
            preSum[curXor] += 1
        return res
