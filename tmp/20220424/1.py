from functools import reduce
from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def intersection(self, nums: List[List[int]]) -> List[int]:
        return sorted(reduce(lambda x, y: x & y, map(set, nums)))
