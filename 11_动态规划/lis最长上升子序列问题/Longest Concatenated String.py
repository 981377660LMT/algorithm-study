# n ≤ 1,000
# 首尾相同的字母可以相连，求最大长度，使得开头和结尾相同
from functools import lru_cache


class Solution:
    def solve(self, words):
        """Return the longest resulting string possible such that the first character and the last character is the same."""

        @lru_cache(None)
        def dfs(index: int, startChar: str, endChar: str) -> int:
            if index == n:
                return 0 if startChar == endChar != '' else -int(1e99)

            # 不选
            res = dfs(index + 1, startChar, endChar)

            # 第一次选
            if startChar == '':
                res = max(
                    res, dfs(index + 1, words[index][0], words[index][-1]) + len(words[index])
                )

            # 非第一次选
            elif words[index][0] == endChar:
                res = max(res, dfs(index + 1, startChar, words[index][-1]) + len(words[index]))

            return res

        n = len(words)
        res = dfs(0, '', '')
        dfs.cache_clear()
        return res if res > 0 else 0


# 18  1 5
print(Solution().solve(words=["hello", "olympic", "crunch"]))
print(Solution().solve(words=["b"]))
print(Solution().solve(words=["ac"]))
