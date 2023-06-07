# !字符串分割/字符串拼接
# 给你一个字符串 s 和一个字符串列表 wordDict 作为字典。
# !请你判断是否可以利用字典中出现的单词拼接出 s 。
# 注意：不要求字典中出现的单词全部都使用，并且字典中的单词可以重复使用。

# 枚举分割点dp
# 可以使用 Trie/哈希 优化查询字符串是否存在于字典中
# n<=300


from functools import lru_cache
from typing import List


class Solution:
    def wordBreak(self, s: str, wordDict: List[str]) -> bool:
        @lru_cache(None)
        def dfs(index: int) -> bool:
            if index >= n:
                return True
            for i in range(index + 1, n + 1):
                cur = s[index:i]
                if cur in ok and dfs(i):
                    return True
            return False

        n = len(s)
        ok = set(wordDict)
        res = dfs(0)
        dfs.cache_clear()
        return res


assert Solution().wordBreak(s="leetcode", wordDict=["leet", "code"])
