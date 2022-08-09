from collections import defaultdict
from functools import lru_cache
import string


class Solution:
    def longestIdealString(self, s: str, k: int) -> int:
        """endswith子序列dp"""
        dp = defaultdict(int)  # 每种字符结尾的子序列最大长度
        for char in s:
            # !如果这里个数多 可以线段树维护区间最大值
            dp[char] = max(
                1 + dp[pre] for pre in string.ascii_lowercase if abs(ord(char) - ord(pre)) <= k
            )
        return max(dp.values())

    def longestIdealString2(self, s: str, k: int) -> int:
        n = len(s)

        @lru_cache(None)
        def dfs(index: int, pre: int) -> int:
            if index == n:
                return 0
            hash = index * 27 + pre
            if memo[hash] != -1:
                return memo[hash]

            if pre == 0:
                cand1 = dfs(index + 1, pre)
                cand2 = dfs(index + 1, ord(s[index]) - 96) + 1
                res = cand1 if cand1 > cand2 else cand2
                memo[hash] = res
                return res
            else:
                res = dfs(index + 1, pre)
                if abs(ord(s[index]) - 96 - pre) <= k:
                    cand = dfs(index + 1, ord(s[index]) - 96) + 1
                    res = cand if cand > res else res
                memo[hash] = res
                return res

        memo = [-1] * (n * 28)
        res = dfs(0, 0)
        return res


print(Solution().longestIdealString(s="pvjcci", k=4))
