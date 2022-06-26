from functools import lru_cache
from typing import List, Tuple
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)

# 6100. 统计放置房子的方式数
# 街道的两侧各有 n 个地块。每一边的地块都按从 1 到 n 编号
# 现要求街道同一侧不能存在两所房子相邻的情况，请你计算并返回放置房屋的方式数目

# !选或者不选
# !从dp[i-1] 和 dp[i-2] 转移过来


dp = [1, 2]
for i in range(int(1e4) + 10):
    dp.append((dp[-1] + dp[-2]) % MOD)


class Solution:
    def countHousePlacements(self, n: int) -> int:
        """场外打表"""
        res = dp[n]
        return (res * res) % MOD

    def countHousePlacements2(self, n: int) -> int:
        """记忆化dfs"""

        @lru_cache(None)
        def dfs(index: int, pre: bool) -> int:
            """不相邻放置"""
            if index == n:
                return 1

            res = dfs(index + 1, False)  # 不放
            if not pre:
                res += dfs(index + 1, True)

            return res % MOD

        res = dfs(0, False)
        dfs.cache_clear()
        return (res * res) % MOD


# 500478595
