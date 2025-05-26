# https://leetcode.cn/problems/lexicographically-smallest-string-after-adjacent-removals/solutions/3685460/zi-fu-xiao-xiao-le-qu-jian-dp-xian-xing-xmaqk/
# 字符消消乐：区间 DP + 线性 DP
# 消除字符满足如下性质：
#
# 可以消除相邻的「连续」字符。
# 相邻字符消除后，原本不相邻的字符会变成相邻，可以继续消除。换句话说，设子串 A=x+B+y，如果 x 和 y 是「连续」字符，且子串 B 可以完全消除，那么子串 A 可以完全消除。
# 设子串 A=B+C，如果子串 B 和 C 可以完全消除，那么子串 A 可以完全消除。

from functools import lru_cache


def canRemove(x: str, y: str) -> bool:
    d = abs(ord(x) - ord(y))
    return d == 1 or d == 25


class Solution:
    def lexicographicallySmallestString(self, s: str) -> str:
        n = len(s)

        @lru_cache(None)
        def canBeEmpty(i: int, j: int) -> bool:
            if i > j:
                return True
            if canRemove(s[i], s[j]) and canBeEmpty(i + 1, j - 1):
                return True
            for k in range(i + 1, j - 1, 2):
                if canBeEmpty(i, k) and canBeEmpty(k + 1, j):
                    return True
            return False

        @lru_cache(None)
        def dfs(i: int) -> str:
            if i == n:
                return ""
            res = s[i] + dfs(i + 1)
            for j in range(i + 1, n, 2):
                if canBeEmpty(i, j):
                    res = min(res, dfs(j + 1))
            return res

        res = dfs(0)
        dfs.cache_clear()
        canBeEmpty.cache_clear()
        return res
