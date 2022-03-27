from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def findDifference(self, nums1: List[int], nums2: List[int]) -> List[List[int]]:
        res = []
        s1, s2 = set(nums1), set(nums2)
        res.append(list(s1 - s2))
        res.append(list(s2 - s1))
        return res

