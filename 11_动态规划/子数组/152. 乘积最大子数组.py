# 152. 乘积最大子数组
# https://leetcode.cn/problems/maximum-product-subarray/


from typing import List


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


assert maxProduct([2, 3, -2, 4]) == 6
