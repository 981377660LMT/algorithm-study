from typing import List
from collections import defaultdict


class Solution:
    def brightestPosition(self, lights: List[List[int]]) -> int:
        diff = defaultdict(int)
        for start, range in lights:
            diff[start - range] += 1
            diff[start + range + 1] -= 1

        res = -1
        curSum = 0
        maxSum = -int(1e20)
        for pos, delta in sorted(diff.items()):
            curSum += delta
            if curSum > maxSum:
                res = pos
                maxSum = curSum
        return res

 


print(Solution().brightestPosition(lights=[[-3, 2], [1, 2], [3, 3]]))
