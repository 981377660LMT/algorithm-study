from functools import cache, lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个 正 整数 n 和 x 。

# 请你返回将 n 表示成一些 互不相同 正整数的 x 次幂之和的方案数。换句话说，你需要返回互不相同整数 [n1, n2, ..., nk] 的集合数目，满足 n = n1x + n2x + ... + nkx 。

# 由于答案可能非常大，请你将它对 109 + 7 取余后返回。


# 比方说，n = 160 且 x = 3 ，一个表示 n 的方法是 n = 23 + 33 + 53 。

pow = lru_cache(None)(pow)


class Solution:
    def numberOfWays(self, n: int, x: int) -> int:
        res = 0

        @lru_cache(None)
        def dfs(base: int, remain: int) -> int:
            if remain <= 0:
                return remain == 0
            if base > n:
                return 0

            pow_ = pow(base, x)
            if pow_ > remain:
                return 0
            res1 = dfs(base + 1, remain - pow_)
            res2 = dfs(base + 1, remain)
            return (res1 + res2) % MOD

        res = dfs(1, n)
        dfs.cache_clear()
        return res % MOD


print(Solution().numberOfWays(4, 2))
