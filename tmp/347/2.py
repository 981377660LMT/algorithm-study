from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始的字符串 s 和一个单词字典 dictionary 。你需要将 s 分割成若干个 互不重叠 的子字符串，每个子字符串都在 dictionary 中出现过。s 中可能会有一些 额外的字符 不在任何子字符串中。

# 请你采取最优策略分割 s ，使剩下的字符 最少 。


# !字符串分割类型dp
class Solution:
    def minExtraChar(self, s: str, dictionary: List[str]) -> int:
        res = len(s)
        ok = set(dictionary)

        @lru_cache(None)
        def dfs(index: int) -> int:
            if index >= len(s):
                return 0
            res = dfs(index + 1)
            for i in range(index + 1, len(s) + 1):
                if s[index:i] in ok:
                    res = max(res, dfs(i) + i - index)
            return res

        return len(s) - dfs(0)


# s = "leetscode", dictionary = ["leet","code","leetcode"]
print(Solution().minExtraChar("leetscode", ["leet", "code", "leetcode"]))
