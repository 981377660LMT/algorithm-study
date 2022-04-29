from itertools import accumulate
import sys
from functools import lru_cache

sys.setrecursionlimit(int(1e9))


# 每行最多k个字符

# 前缀和+线性dp
class Solution:
    def minimumCost(self, sentence: str, k: int) -> int:
        @lru_cache(None)
        def dfs(pos: int) -> int:
            """dfs(pos)表示从第pos个单词开始分隔的最小成本(1<=pos<=n)
            
            总成本就是除开最后一行以外的其它所有行的分隔成本之和。
            """
            if calLength(pos, n) <= k:  # 最后一行
                return 0

            res = int(1e20)
            for nPos in range(pos, n + 1):
                length = calLength(pos, nPos)
                if length > k:
                    break
                cost = (k - length) * (k - length)
                res = min(res, dfs(nPos + 1) + cost)
            return res

        def calLength(start: int, end: int) -> int:
            """第start个单词到第end个单词变为1行的长度"""
            assert 1 <= start <= end <= n
            return preSum[end] - preSum[start - 1] + (end - start)

        lens = list(map(len, sentence.split(' ')))
        n = len(lens)
        preSum = list(accumulate(lens, initial=0))
        res = dfs(1)
        dfs.cache_clear()
        return res

    # def minimumCost(self, sentence: str, k: int) -> int:
    #     """搜索的时候也可以剪枝"""

    #     @lru_cache(None)
    #     def dfs(pos: int, width: int, pathSum: int) -> None:
    #         nonlocal res

    #         # 剪枝1:超出res就返回
    #         if pathSum > res:
    #             return

    #         if index == len(lens):
    #             res = min(res, pathSum)
    #             return

    #         # 剪枝2:优先把好的放前面搜素
    #         if width + lens[index] + 1 <= k:
    #             dfs(index + 1, width + lens[index] + 1, pathSum)
    #         dfs(index + 1, lens[index], pathSum + (k - width) ** 2)

    #     lens = list(map(len, sentence.split(' ')))
    #     res = 0x7FFFFFFF

    #     dfs(1, lens[0], 0)
    #     return res


print(Solution().minimumCost(sentence="i love leetcode", k=12))
