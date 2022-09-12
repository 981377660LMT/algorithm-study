from itertools import accumulate
from typing import List


class Solution:
    def meetRequirement(self, n: int, lights: List[List[int]], requirement: List[int]) -> int:
        diff = [0] * (n + 1)
        for start, range in lights:
            diff[max(0, start - range)] += 1
            diff[min(n, start + range + 1)] -= 1
        diff = list(accumulate(diff))
        return sum(d >= r for d, r in zip(diff, requirement))
