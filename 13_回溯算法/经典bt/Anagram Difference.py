# 使两个anagram相等的最小交换次数
from typing import List

# n<=25


class Solution:
    def solve(self, s0, s1):
        def bt(i: int, nums1: List[int], nums2: List[int]) -> int:
            if i == n:
                return 0
            if nums1[i] == nums2[i]:
                return bt(i + 1, nums1, nums2)

            res = n
            for j in range(i, n):
                if nums1[j] == nums2[i]:
                    nums1[i], nums1[j] = nums1[j], nums1[i]
                    res = min(res, 1 + bt(i + 1, nums1, nums2))
                    nums1[i], nums1[j] = nums1[j], nums1[i]
            return res

        n = len(s0)
        return bt(0, list(s0), list(s1))


print(Solution().solve('dom', 'mod'))

