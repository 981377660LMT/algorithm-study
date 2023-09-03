from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的字符串 num ，表示一个非负整数。

# 在一次操作中，您可以选择 num 的任意一位数字并将其删除。请注意，如果你删除 num 中的所有数字，则 num 变为 0。

# 返回最少需要多少次操作可以使 num 变成特殊数字。


# 如果整数 x 能被 25 整除，则该整数 x 被认为是特殊数字
class Solution:
    def minimumOperations(self, num: str) -> int:
        @lru_cache(None)
        def dfs(index: int, mod: int) -> int:
            if index == len(num):
                return 0 if mod == 0 else -INF
            res1 = dfs(index + 1, mod)
            res2 = dfs(index + 1, (mod * 10 + int(num[index])) % 25) + 1
            return max(res1, res2)

        res = dfs(0, 0)
        dfs.cache_clear()
        return len(num) - res
