# 对于范围 [0, n - 1] 的每个下标 i ，你可以将 nums1[i] 或 nums2[i] 的值赋给 nums3[i] 。

# 你的任务是使用最优策略为 nums3 赋值，以最大化 nums3 中 最长非递减子数组 的长度。


from typing import List
from LIS模板 import LIS


class Solution:
    def maxNonDecreasingLength(self, nums1: List[int], nums2: List[int]) -> int:
        """两个数组选数,使得最长非递减`子数组`最长."""
        n, res = len(nums1), 1
        dp = [1, 1]
        for i in range(1, n):
            ndp = [1, 1]
            for pre in range(2):
                preNums = nums1 if pre == 0 else nums2
                for cur in range(2):
                    curNums = nums1 if cur == 0 else nums2
                    if preNums[i - 1] <= curNums[i]:
                        ndp[cur] = max(ndp[cur], dp[pre] + 1)
            res = max(res, ndp[0], ndp[1])
            dp = ndp
        return res

    def maxNonDecreasingLength2(self, nums1: List[int], nums2: List[int]) -> int:
        """两个数组选数,使得最长非递减`子序列`最长->LIS."""
        nums = []
        for a, b in zip(nums1, nums2):
            if a == b:
                nums.append(a)
            else:
                nums.extend(sorted([a, b], reverse=True))
        return LIS(nums, isStrict=False)
