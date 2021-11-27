from typing import List
from collections import defaultdict

# 我们可以在 word1 的任何地方添加一个字母使其变成 word2
# 考虑每个以当前word结尾的串 dp
class Solution:
    def longestStrChain(self, words: List[str]) -> int:
        res = 0
        dp = defaultdict(int)
        for w in sorted(words, key=len):
            for remove in range(len(w)):
                dp[w] = max(dp[w], dp[w[:remove] + w[(remove + 1) :]] + 1)
            res = max(res, dp[w])
        return res


assert Solution().longestStrChain(["a", "b", "ba", "bca", "bda", "bdca"]) == 4

# 输出：4
# 解释：最长单词链之一为 "a","ba","bda","bdca"。
