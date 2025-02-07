# 3434. 子数组操作后的最大频率
# https://leetcode.cn/problems/maximum-frequency-after-subarray-operation/description/
# 给你一个长度为 n 的数组 nums ，同时给你一个整数 k 。
# 你可以对 nums 执行以下操作 一次 ：
# 选择一个子数组 nums[i..j] ，其中 0 <= i <= j <= n - 1 。
# 选择一个整数 x 并将 nums[i..j] 中 所有 元素都增加 x 。
# 请你返回执行以上操作以后数组中 k 出现的 最大 频率。
# 子数组 是一个数组中一段连续 非空 的元素序列。
#
# 1 <= n == nums.length <= 105
# 1 <= nums[i] <= 50
# 1 <= k <= 50
#
# !枚举+状态机dp

from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxFrequency(self, nums: List[int], k: int) -> int:
        f0 = max_f12 = 0
        f1 = [0] * 51
        for x in nums:
            if x == k:
                max_f12 += 1
                f0 += 1
            else:
                f1[x] = max2(f1[x], f0) + 1
                max_f12 = max2(max_f12, f1[x])
        return max_f12
