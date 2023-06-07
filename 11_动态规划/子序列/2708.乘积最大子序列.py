# https://leetcode.cn/problems/maximum-product-subarray/
# !乘积最大子序列
# 给你一个下标从 0 开始的整数数组 nums ，它表示一个班级中所有学生在一次考试中的成绩。
# 老师想选出一部分同学组成一个 非空 小组，且这个小组的 实力值 最大，
# 如果这个小组里的学生下标为 i0, i1, i2, ... , ik ，
# 那么这个小组的实力值定义为 nums[i0] * nums[i1] * nums[i2] * ... * nums[ik​] 。

# 请你返回老师创建的小组能得到的最大实力值为多少。
# n<=1e5

# !和乘积最大子数组相比,多了一个`不取`的选项.

from typing import List


def maxStrength(nums: List[int]) -> int:
    """
    乘积最大子序列.
    https://leetcode.cn/problems/maximum-strength-of-a-group/
    """
    res, dpMin, dpMax = nums[0], nums[0], nums[0]
    for v in nums[1:]:
        ndpMin = min(dpMin * v, dpMax * v, v, dpMin)
        ndpMax = max(dpMin * v, dpMax * v, v, dpMax)
        dpMin, dpMax = ndpMin, ndpMax
        res = max(res, dpMax)
    return res


def maxProduct(nums: List[int]) -> int:
    """
    乘积最大子数组.
    https://leetcode.cn/problems/maximum-product-subarray/
    """
    res, dpMin, dpMax = nums[0], nums[0], nums[0]
    for v in nums[1:]:
        ndpMin = min(dpMin * v, dpMax * v, v)
        ndpMax = max(dpMin * v, dpMax * v, v)
        dpMin, dpMax = ndpMin, ndpMax
        res = max(res, dpMax)
    return res
