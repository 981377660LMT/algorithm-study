# 1534. 统计好三元组


from itertools import combinations
from typing import List


class Solution:
    def countGoodTriplets(self, arr: List[int], a: int, b: int, c: int) -> int:
        return sum(
            abs(i - j) <= a and abs(j - k) <= b and abs(i - k) <= c
            for i, j, k in combinations(arr, 3)
        )
