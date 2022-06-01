from functools import lru_cache
from typing import List
from collections import defaultdict

# 我们可以在 word1 的任何地方添加一个字母使其变成 word2
# 考虑每个以当前word结尾的串 dp
# 1 <= words.length <= 1000
# 1 <= words[i].length <= 16
# words[i] 仅由小写英文字母组成。


class Solution:
    def longestStrChain(self, words: List[str]) -> int:
        res = 0
        dp = defaultdict(int)  # 每个单词结尾的最长链长度

        for cur in sorted(words, key=len):
            for i in range(len(cur)):
                pre = cur[:i] + cur[i + 1 :]
                dp[cur] = max(dp[cur], dp[pre] + 1)
            res = max(res, dp[cur])

        return res

    def longestStrChain2(self, words: List[str]) -> int:
        """建图求最长路  不需要邻接表建图 只需要genNexts函数找临边 (广义邻居)"""

        def genNexts(cur: str):
            yield from (cur[:i] + cur[i + 1 :] for i in range(len(cur)))

        @lru_cache(None)
        def dfs(cur: str) -> int:
            res = 1
            for next in genNexts(cur):
                if next in words:
                    res = max(res, dfs(next) + 1)
            return res

        words = set(words)
        return max(dfs(w) for w in words)


assert Solution().longestStrChain(["a", "b", "ba", "bca", "bda", "bdca"]) == 4

# 输出：4
# 解释：最长单词链之一为 "a","ba","bda","bdca"。
