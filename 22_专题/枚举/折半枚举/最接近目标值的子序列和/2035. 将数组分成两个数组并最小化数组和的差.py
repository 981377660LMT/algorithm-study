# 给你一个长度为 2 * n 的整数数组。你需要将 nums 分成 两个 长度为 n 的数组，
# 分别求出两个数组的和，
# 并 最小化 两个数组和之 差的绝对值 。nums 中每个元素都需要放入两个数组之一。


from collections import defaultdict
from typing import DefaultDict, List, Set


def subsetSum(nums: List[int]) -> DefaultDict[int, Set[int]]:
    """O(2^n)求所有子集的可能和"""
    n = len(nums)
    res = defaultdict(set, {0: set([0])})
    dp = [0] * (1 << n)
    for i in range(n):
        for pre in range(1 << i):
            cur = pre + (1 << i)
            dp[cur] = dp[pre] + nums[i]
            res[cur.bit_count()].add(dp[cur])
    return res


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
    def minimumDifference(self, nums: List[int]) -> int:
        """x+y=sum 求abs(x-y)最小值  即abs(2*x-sum)最小值"""
        n = len(nums)
        target = sum(nums)
        nums = [num * 2 for num in nums]
        sums1, sums2 = (
            subsetSum(nums[: n // 2]),
            subsetSum(nums[n // 2 :]),
        )

        res = int(1e20)
        for leftCount in range(n // 2 + 1):
            left, right = sorted(sums1[leftCount]), sorted(sums2[n // 2 - leftCount])
            res = min(res, twoSum(left, right, target))
        return res


print(Solution().minimumDifference([4, 2, 1, 3]))
print(Solution().minimumDifference([3, 9, 7, 3]))

