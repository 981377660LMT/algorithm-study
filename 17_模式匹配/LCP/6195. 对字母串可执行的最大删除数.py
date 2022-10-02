from functools import lru_cache
from getLCP import getLCP


# 给你一个仅由小写英文字母组成的字符串 s 。在一步操作中，你可以：

# 删除 整个字符串 s ，或者
# 对于满足 1 <= i <= s.length / 2 的任意 i ，如果 s 中的 前 i 个字母和接下来的 i 个字母 相等 ，
# 删除 前 i 个字母。
# 例如，如果 s = "ababc" ，那么在一步操作中，你可以删除 s 的前两个字母得到 "abc" ，
# 因为 s 的前两个字母和接下来的两个字母都等于 "ab" 。

# !返回删除 s 所需的最大操作数。
# n<=4000

# !删哪里?枚举所有位置dp 本质上是暴力


# 教训:
# !不要用字符串当记忆化dfs参数 要用index当记忆化参数


class Solution:
    def deleteString(self, s: str) -> int:
        @lru_cache(None)
        def dfs(index: int) -> int:
            if index == n:
                return 0

            remain = n - index
            res = 1
            for i in range(1, remain // 2 + 1):
                if s[index : i + index] == s[i + index : i + i + index]:
                    # if lcp[index][i + index] >= i:  LCP没有暴力比较快
                    cand = dfs(i + index) + 1
                    res = cand if cand > res else res
            return res

        n = len(s)
        # lcp = getLCP1(s)
        res = dfs(0)
        dfs.cache_clear()
        return res


print(Solution().deleteString(s="abcabcdabc"))
print(Solution().deleteString(s="aaabaab"))
print(Solution().deleteString(s="aaaaa"))
