from bisect import bisect_right
from typing import List

MOD = int(1e9 + 7)


class Solution:
    def minAbsoluteSumDiff(self, nums1: List[int], nums2: List[int]) -> int:
        sl = sorted(nums1)
        res = 0
        max_ = -int(1e20)
        n = len(nums1)
        for i in range(n):
            res += abs(nums1[i] - nums2[i])

            # 左、右cand bisect_right更方便
            pos = bisect_right(sl, nums2[i])
            if pos < n:
                max_ = max(max_, abs(nums1[i] - nums2[i]) - abs(sl[pos] - nums2[i]))
            if pos - 1 >= 0:
                max_ = max(max_, abs(nums1[i] - nums2[i]) - abs(sl[pos - 1] - nums2[i]))

        return (res - max_) % MOD

