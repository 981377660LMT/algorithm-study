from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def findClosestNumber(self, nums: List[int]) -> int:
        res = 0
        diff = int(1e20)

        for num in nums:
            delta = abs(num)
            if delta < diff:
                res = num
                diff = delta
            elif delta == diff and num > res:
                res = num

        return res

