import math
from heapq import nlargest, nsmallest


class Solution:
    def solve(self, nums):
        return max(math.prod(nsmallest(2, nums)), math.prod(nlargest(2, nums)))
