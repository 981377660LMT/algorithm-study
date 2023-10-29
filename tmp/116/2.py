from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个长度为偶数下标从 0 开始的二进制字符串 s 。

# 如果可以将一个字符串分割成一个或者更多满足以下条件的子字符串，那么我们称这个字符串是 美丽的 ：

# 每个子字符串的长度都是 偶数 。
# 每个子字符串都 只 包含 1 或 只 包含 0 。
# 你可以将 s 中任一字符改成 0 或者 1 。

# 请你返回让字符串 s 美丽的 最少 字符修改次数。


def min(a, b):
    return a if a < b else b


class Solution:
    def minChanges(self, s: str) -> int:
        @lru_cache(None)
        def dfs(index: int, lenMod: int, zeroOne: int) -> int:
            if index == n:
                return 0 if lenMod == 0 else INF

            # 在此处分割
            res = INF
            cur = nums[index]
            if lenMod == 0:
                res = min(res, dfs(index + 1, 1, cur))

            # 不分割,沿用上一次的分割
            res = min(res, dfs(index + 1, lenMod ^ 1, zeroOne) + (1 if cur != zeroOne else 0))
            return res

        n = len(s)
        nums = [int(i) for i in s]
        res = dfs(1, 1, nums[0])
        dfs.cache_clear()
        return res
