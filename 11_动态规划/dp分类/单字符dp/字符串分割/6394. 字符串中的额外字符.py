# https://leetcode.cn/problems/extra-characters-in-a-string/
# n<=50

# 给你一个下标从 0 开始的字符串 s 和一个单词字典 dictionary 。
# 你需要将 s 分割成若干个 互不重叠 的子字符串，
# 每个子字符串都在 dictionary 中出现过。
# s 中可能会有一些 额外的字符 不在任何子字符串中。

# 请你采取最优策略分割 s ，使剩下的字符 最少 。


from functools import lru_cache
from typing import List

INF = int(1e18)


class Solution:
    def minExtraChar(self, s: str, dictionary: List[str]) -> int:
        @lru_cache(None)
        def dfs(index: int) -> int:
            """返回从index开始的最大匹配个数."""
            if index == n:
                return 0

            res = dfs(index + 1)  # !jump
            for i in range(index + 1, n + 1):  # !not jump
                cur = s[index:i]
                if cur in ok:
                    res = max(res, len(cur) + dfs(i))
            return res

        n = len(s)
        ok = set(dictionary)
        res = dfs(0)
        dfs.cache_clear()
        return len(s) - res

    def minExtraChar2(self, s: str, dictionary: List[str]) -> int:
        n = len(s)
        ok = set(dictionary)
        dp = [0] * (n + 1)  # dp[i]表示前i个字符的最大匹配个数
        for i in range(1, n + 1):
            dp[i] = dp[i - 1]  # !jump
            for j in range(i):  # !not jump
                cur = s[j:i]
                if cur in ok:
                    dp[i] = max(dp[i], len(cur) + dp[j])
        return len(s) - dp[-1]

    def minExtraChar3(self, s: str, dictionary: List[str]) -> int:
        """预处理转移点."""
        n = len(s)
        pre = [[] for _ in range(n)]  # !每个索引可以从哪些索引转移过来
        dp = [0] * (n + 1)

        for word in dictionary:  # 这一段可以AC自动机优化
            len_ = len(word) - 1
            start = s.find(word)
            while start != -1:
                pre[start + len_].append(start)
                start = s.find(word, start + 1)

        for i, pos in enumerate(pre):
            if not pos:
                dp[i + 1] = dp[i] + 1
                continue

            res = dp[i] + 1
            for p in pos:
                res = min(res, dp[p])
            dp[i + 1] = res

        return dp[-1]
