import functools
import sys
from functools import lru_cache

sys.setrecursionlimit(1000000)

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
        def dfs(index, width) -> int:
            if index == n:
                return 0
            res = dfs(index + 1, word_len[index]) + (k - width) ** 2
            if width + word_len[index] + 1 <= k:
                res = min(res, dfs(index + 1, width + word_len[index] + 1))
            return res

        return dfs(1, word_len[0])

    # 改用dfs+剪枝
    def minimumCost(self, sentence: str, k: int) -> int:

        words = sentence.split(' ')
        n = len(words)
        word_len = []
        for w in words:
            word_len.append(len(w))

        res = 0x7FFFFFFF

        @functools.lru_cache(None)
        def dfs(index: int, width: int, pathSum: int) -> None:
            nonlocal res

            # 关键的剪枝1:超出res就返回
            if pathSum > res:
                return

            if index == n:
                res = min(res, pathSum)
                return

            # 2.这是更最关键的剪枝:优先把好的放前面搜素
            if width + word_len[index] + 1 <= k:
                dfs(index + 1, width + word_len[index] + 1, pathSum)
            dfs(index + 1, word_len[index], pathSum + (k - width) ** 2)

        dfs(1, word_len[0], 0)
        return res


print(Solution().minimumCost(sentence="i love leetcode", k=12))
