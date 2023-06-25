from typing import List
from sortedcontainers import SortedList


class Solution:
    def goodTriplets(self, nums1: List[int], nums2: List[int]) -> int:
        n = len(nums1)
        indexByValue = dict({value: index for index, value in enumerate(nums1)})
        target = [indexByValue[num] for num in nums2]

        res = 0
        sl = SortedList()

        for i, num in enumerate(target):
            sl.add(num)
            smaller = sl.bisect_left(num)
            bigger = n - 1 - num - (i - smaller)
            res += smaller * bigger
        return res


print(Solution().goodTriplets(nums1=[2, 0, 1, 3], nums2=[0, 1, 2, 3]))
print(Solution().goodTriplets(nums1=[4, 0, 1, 3, 2], nums2=[4, 1, 0, 2, 3]))
