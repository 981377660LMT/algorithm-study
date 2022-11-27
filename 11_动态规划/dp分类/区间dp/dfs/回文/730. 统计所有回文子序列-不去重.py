from functools import lru_cache


# 给定一个字符串 s，返回 s 中不同的非空「回文子序列」个数 。
# !字符串 S 的长度将在[1, 1000]范围内。
# 每个字符 S[i] 将会是集合 {'a', 'b', 'c', 'd'} 中的某一个。

MOD = int(1e9 + 7)


class Solution:
    def countPalindromicSubsequences3(self, s: str) -> int:
        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            """
            [l,r] 间的回文子序列数
            1. 如果两个端点不全取:
            `dfs(l, r - 1) + dfs(l + 1, r) - dfs(l + 1, r - 1)`
            2. 如果取两个端点:
            那么就要加上 `dfs(l + 1, r - 1) + 1`
            """
            if left >= right:
                return int(left == right)
            if s[left] == s[right]:
                return (1 + dfs(left + 1, right) + dfs(left, right - 1)) % MOD
            return (dfs(left, right - 1) + dfs(left + 1, right) - dfs(left + 1, right - 1)) % MOD

        return dfs(0, len(s) - 1)


print(Solution().countPalindromicSubsequences3("abcd"))  # 4
print(Solution().countPalindromicSubsequences3("aab"))  # 4
print(Solution().countPalindromicSubsequences3("aaaa"))  # 15
