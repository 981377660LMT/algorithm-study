# !字符串分割/字符串拼接
# 给你一个字符串 s 和一个字符串列表 wordDict 作为字典。
# !请你判断是否可以利用字典中出现的单词拼接出 s 。
# 注意：不要求字典中出现的单词全部都使用，并且字典中的单词可以重复使用。

# 枚举分割点dp
# 可以使用 Trie/哈希 优化查询字符串是否存在于字典中
# n<=300


from functools import lru_cache
from typing import List

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def wordBreak(self, s: str, wordDict: List[str]) -> bool:
        trie = dict()
        for word in wordDict:
            root = trie
            for c in word:
                if c not in root:
                    root[c] = dict()
                root = root[c]
            root["X"] = True
        n = len(s)
        dp = [INF] * (n + 1)
        dp[0] = True
        ptr = 0
        while ptr < n + 1:
            if dp[ptr] != INF:
                j = ptr + 1
                root = trie
                while j < n + 1:
                    if s[j - 1] in root:
                        root = root[s[j - 1]]
                        if "X" in root:
                            dp[j] = True
                        j += 1
                    else:
                        break
            ptr += 1
        return dp[-1]


assert Solution().wordBreak(s="leetcode", wordDict=["leet", "code"])
