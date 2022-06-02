# 给你一个整数数组 nums 和一个目标值 goal 。
# 你需要从 nums 中选出一个子序列(子集,（可能全部或无）)，使子序列元素总和最接近 goal
# 1 <= nums.length <= 40
# -107 <= nums[i] <= 107


from typing import List, Set


def subsetSum(nums: List[int]) -> Set[int]:
    """O(2^n)求所有子集的可能和"""
    dp = set({0})
    for cur in nums:
        dp |= {(cur + pre) for pre in (dp | {0})}
    return dp


def twoSum(nums1: List[int], nums2: List[int], target: int) -> int:
    i, j = 0, len(nums2) - 1
    res = int(1e20)
    while i < len(nums1) and j >= 0:
        cand = nums1[i] + nums2[j]
        res = min(res, abs(cand - target))
        if cand < target:
            i += 1
        elif cand > target:
            j -= 1
        else:
            return 0
    return res


class Solution:
    def minAbsDifference(self, nums: List[int], goal: int) -> int:
        n = len(nums)
        sums1, sums2 = sorted(subsetSum(nums[: n // 2])), sorted(subsetSum(nums[n // 2 :]))
        return twoSum(sums1, sums2, goal)

