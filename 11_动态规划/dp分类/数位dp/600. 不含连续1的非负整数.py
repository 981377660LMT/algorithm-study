from functools import lru_cache


@lru_cache(None)
def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, pre: int, isLimit: bool) -> int:
        """当前在第pos位，前一位为pre，isLimit表示是否贴合上界"""
        if pos == 0:
            return 1

        res = 0
        up = nums[pos - 1] if isLimit else 1
        for cur in range(up + 1):
            if pre == cur == 1:
                continue
            res += dfs(pos - 1, cur, (isLimit and cur == up))
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, 2)
        nums.append(mod)
        upper = div
    return dfs(len(nums), -1, True)


class Solution:
    def findIntegers(self, n: int) -> int:
        """给定一个正整数 n ，返回范围在 [0, n] 都非负整数中，其二进制表示不包含 连续的 1 的个数。"""
        return cal(n)

