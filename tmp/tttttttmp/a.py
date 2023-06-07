from heapq import nlargest, nsmallest
from typing import List


class Solution:
    def buyChoco(self, prices: List[int], money: int) -> int:
        minSum = sum(nsmallest(2, prices))
        if money < minSum:
            return money
        return money - minSum
