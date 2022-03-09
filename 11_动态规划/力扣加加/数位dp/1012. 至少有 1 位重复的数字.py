from functools import lru_cache

# 1 <= n <= 109


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, visited: int, isOk: bool, hasLeadingZero: bool, isLimit: bool) -> int:
        """当前在第pos位，visited表示每个数使用的状态，isOk表示是否有重复数字，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == 0:
            return isOk

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos - 1, visited, isOk, True, (isLimit and cur == up))
            else:
                res += dfs(
                    pos - 1,
                    (visited | (1 << cur)),
                    not not (isOk or (visited & (1 << cur))),
                    False,
                    (isLimit and cur == up),
                )
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, 10)
        nums.append(mod)
        upper = div
    res = dfs(len(nums), 0, False, True, True, '')
    dfs.cache_clear()
    return res


# 给定正整数 n，返回在 [1, n] 范围内具有 至少 1 位 重复数字的正整数的个数。
class Solution:
    def numDupDigitsAtMostN(self, n: int) -> int:
        """给定正整数 n，返回在 [1, n] 范围内具有 至少 1 位 重复数字的正整数的个数。"""
        return cal(n) - cal(0)


print(Solution().numDupDigitsAtMostN(20))
