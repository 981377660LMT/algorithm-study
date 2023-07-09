# 给你一个二进制字符串 s ，你需要将字符串分割成一个或者多个子字符串，使每个子字符串都是美丽的。
# 如果一个字符串满足以下条件，我们称它是 美丽 的：
# !1.它不包含前导 0 。
# !2.它是 5 的幂的 二进制 表示。
# 请你返回分割后的子字符串的 最少 数目。如果无法将字符串 s 分割成美丽子字符串，请你返回 -1 。

from functools import lru_cache

INF = int(1e18)
POW5 = [bin(5**i)[2:] for i in range(20)]


class Solution:
    def minimumBeautifulSubstrings(self, s: str) -> int:
        n = len(s)

        @lru_cache(None)
        def dfs(index: int) -> int:
            if index == n:
                return 0
            if s[index] == "0":
                return INF

            res = INF
            for num in POW5:
                if index + len(num) > n:
                    break
                if s[index : index + len(num)] == num:
                    res = min(res, dfs(index + len(num)) + 1)
            return res

        n = len(s)
        res = dfs(0)
        dfs.cache_clear()
        return res if res < INF else -1
