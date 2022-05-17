from collections import defaultdict
from typing import List

# 1 <= A.length <= 1000
class Solution:
    def countTriplets(self, A: List[int]) -> int:
        memo = defaultdict(int)
        for n1 in A:
            for n2 in A:
                memo[n1 & n2] += 1

        res = 0
        for num in A:
            for key, val in memo.items():
                if num & key == 0:
                    res += val
        return res
