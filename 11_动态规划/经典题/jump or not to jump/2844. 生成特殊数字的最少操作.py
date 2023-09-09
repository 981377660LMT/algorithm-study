# 2844. 生成特殊数字的最少操作
# https://leetcode.cn/problems/minimum-operations-to-make-a-special-number/
# 给你一个下标从 0 开始的字符串 num ，表示一个非负整数。
# 在一次操作中，您可以选择 num 的任意一位数字并将其删除。请注意，如果你删除 num 中的所有数字，则 num 变为 0。
# 返回最少需要多少次操作可以使 num 变成特殊数字。
# 如果整数 x 能被 25 整除，则该整数 x 被认为是特殊数字

from functools import lru_cache


INF = int(1e20)


class Solution:
    def minimumOperations(self, num: str) -> int:
        @lru_cache(None)
        def dfs(index: int, mod: int) -> int:
            if index == len(num):
                return 0 if mod == 0 else -INF
            res1 = dfs(index + 1, mod)
            res2 = dfs(index + 1, (mod * 10 + int(num[index])) % 25) + 1
            return res1 if res1 > res2 else res2

        res = dfs(0, 0)
        dfs.cache_clear()
        return len(num) - res


if __name__ == "__main__":
    assert Solution().minimumOperations("2245047") == 2
