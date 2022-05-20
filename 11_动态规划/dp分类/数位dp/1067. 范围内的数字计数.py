from functools import lru_cache


@lru_cache(None)
def cal(upper: int, queryDigit: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, count: int) -> int:
        """当前在第pos位，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == len(nums):
            return count

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), count)
            else:
                res += dfs(pos + 1, False, (isLimit and cur == up), count + (cur == queryDigit))
        return res

    nums = list(map(int, str(upper)))
    return dfs(0, True, True, 0)


class Solution:
    def digitsCount(self, d: int, low: int, high: int) -> int:
        # 1 <= low <= high <= 2×10^8
        return cal(high, d) - cal(low - 1, d)
