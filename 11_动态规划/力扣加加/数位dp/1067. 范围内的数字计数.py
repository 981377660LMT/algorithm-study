from functools import lru_cache


@lru_cache(None)
def cal(upper: int, queryDigit: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, count: int, hasLeadingZero: bool, isLimit: bool) -> int:
        """当前在第pos位，出现次数为count，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == 0:
            return count

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos - 1, count, True, (isLimit and cur == up))
            else:
                res += dfs(pos - 1, count + (cur == queryDigit), False, (isLimit and cur == up))
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, 10)
        nums.append(mod)
        upper = div
    return dfs(len(nums), 0, True, True)


class Solution:
    def digitsCount(self, d: int, low: int, high: int) -> int:
        return cal(high, d) - cal(low - 1, d)
