import sys
from functools import lru_cache

sys.setrecursionlimit(int(1e9))

# 有点像放📕那道题
# 1105. 填充书架..py

# 每行最多k个字符


class Solution:
    def minimumCost2(self, sentence: str, k: int) -> int:
        # dp[i] 表示以第i个单词为某行最后一个单词的最小代价
        words = sentence.split(' ')
        n = len(words)
        word_len = []
        for w in words:
            word_len.append(len(w))

        # 为什么记忆化没用了  => 因为 width 不能做状态
        @lru_cache(None)
        def dfs(index: int, width: int) -> int:
            if index == n:
                return 0
            res = dfs(index + 1, word_len[index]) + (k - width) ** 2
            if width + word_len[index] + 1 <= k:
                res = min(res, dfs(index + 1, width + word_len[index] + 1))
            return res

        return dfs(1, word_len[0])

    def minimumCost(self, sentence: str, k: int) -> int:
        """搜索的时候也可以剪枝"""
        word_lens = list(map(len, sentence.split(' ')))
        res = 0x7FFFFFFF

        @lru_cache(None)
        def dfs(index: int, width: int, pathSum: int) -> None:
            nonlocal res

            # 剪枝1:超出res就返回
            if pathSum > res:
                return

            if index == len(word_lens):
                res = min(res, pathSum)
                return

            # 剪枝2:优先把好的放前面搜素
            if width + word_lens[index] + 1 <= k:
                dfs(index + 1, width + word_lens[index] + 1, pathSum)
            dfs(index + 1, word_lens[index], pathSum + (k - width) ** 2)

        dfs(1, word_lens[0], 0)
        return res


print(Solution().minimumCost(sentence="i love leetcode", k=12))
