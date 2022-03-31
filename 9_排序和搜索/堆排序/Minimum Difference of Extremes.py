# 求最值index的最小差值，你有三次将任何数变为任意值的机会
from heapq import nlargest, nsmallest


class Solution:
    def solve(self, nums):
        if len(nums) <= 3:
            return 0

        (m1, m2, m3, m4), (n4, n3, n2, n1) = nsmallest(4, nums), nlargest(4, nums)
        return min(n1 - m1, n2 - m2, n3 - m3, n4 - m4)

