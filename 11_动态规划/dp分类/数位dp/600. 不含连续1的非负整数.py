from functools import lru_cache

# 数位dp模板
@lru_cache(None)
def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, pre: int) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == len(nums):
            return not hasLeadingZero

        res = 0
        up = nums[pos] if isLimit else 1
        for cur in range(up + 1):
            if pre == cur == 1:
                continue
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), cur)
            else:
                res += dfs(pos + 1, False, (isLimit and cur == up), cur)
        return res

    nums = list(map(int, bin(upper)[2:]))
    return dfs(0, True, True, -1)


class Solution:
    def findIntegers(self, n: int) -> int:
        """给定一个正整数 n ，返回范围在 [0, n] 都非负整数中，其二进制表示不包含 连续的 1 的个数。"""
        return cal(n) - cal(0) + 1

