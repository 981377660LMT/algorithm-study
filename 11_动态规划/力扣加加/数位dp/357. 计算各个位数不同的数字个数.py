from functools import lru_cache


@lru_cache(None)
def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, visited: int, hasLeadingZero: bool, isLimit: bool) -> int:
        """当前在第pos位，已经用了visited，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == 0:
            return 1

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            if visited & (1 << cur):
                continue
            if hasLeadingZero and cur == 0:
                res += dfs(pos - 1, visited, True, (isLimit and cur == up))
            else:
                res += dfs(pos - 1, visited | (1 << cur), False, (isLimit and cur == up))
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, 10)
        nums.append(mod)
        upper = div
    return dfs(len(nums), 0, True, True)


class Solution:
    def countNumbersWithUniqueDigits(self, n: int) -> int:
        if n >= 11:
            return 0
        return cal(10 ** n - 1)

