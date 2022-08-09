from functools import lru_cache


@lru_cache(None)
def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, visited: int) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == len(nums):
            return not hasLeadingZero

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if visited & (1 << cur):
                continue
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), visited)
            else:
                res += dfs(pos + 1, False, (isLimit and cur == up), visited | (1 << cur))
        return res

    nums = list(map(int, str(upper)))
    return dfs(0, True, True, 0)


class Solution:
    def countNumbersWithUniqueDigits(self, n: int) -> int:
        return cal(10**n - 1) - cal(0) + 1  # cal(10 ** n - 1) - cal(0) 表示[1,x] 再加上x=0的情况
