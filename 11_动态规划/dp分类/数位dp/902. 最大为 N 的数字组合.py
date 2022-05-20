from functools import lru_cache
from typing import List


def cal(upper: int, digits: List[int]) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: int, isLimit: bool) -> int:
        """当前在第pos位，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == len(nums):
            return not hasLeadingZero

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up))
            else:
                if cur not in digits:
                    continue
                res += dfs(pos + 1, False, (isLimit and cur == up))
        return res

    nums = list(map(int, str(upper)))
    return dfs(0, True, True)


class Solution:
    def atMostNGivenDigitSet(self, digits: List[str], n: int) -> int:
        """返回 可以生成的小于或等于给定整数 n 的正整数的个数 。"""
        return cal(n, list(map(int, digits)))


print(Solution().atMostNGivenDigitSet(digits=["1", "3", "5", "7"], n=100))
