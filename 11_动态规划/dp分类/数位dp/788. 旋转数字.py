from functools import lru_cache
from typing import Set


def cal(upper: int, good: Set[int]) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, isGood: bool) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界,isGood表示是否包含good"""
        if pos == len(nums):
            return isGood and not hasLeadingZero

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), False)
            else:
                if cur in good:
                    res += dfs(
                        pos + 1, False, (isLimit and cur == up), isGood or cur in (2, 5, 6, 9)
                    )
        return res

    nums = list(map(int, str(upper)))
    return dfs(0, True, True, False)


class Solution:
    def rotatedDigits(self, n: int) -> int:
        """[1,n]中有多少个数是好数(可选(0,1,8,2,5,6,9)并且至少有一个(2,5,6,9))"""
        return cal(n, set([0, 1, 8, 2, 5, 6, 9])) - cal(0, set([0, 1, 8, 2, 5, 6, 9]))
