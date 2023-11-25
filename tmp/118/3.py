from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 你在一个水果超市里，货架上摆满了玲琅满目的奇珍异果。

# 给你一个下标从 1 开始的数组 prices ，其中 prices[i] 表示你购买第 i 个水果需要花费的金币数目。

# 水果超市有如下促销活动：

# 如果你花费 price[i] 购买了水果 i ，那么接下来的 i 个水果你都可以免费获得。
# 注意 ，即使你 可以 免费获得水果 j ，你仍然可以花费 prices[j] 个金币去购买它以便能免费获得接下来的 j 个水果。


# 请你返回获得所有水果所需要的 最少 金币数。
def max2(a, b):
    return a if a > b else b


def min2(a, b):
    return a if a < b else b


class Solution:
    def minimumCoins(self, prices: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, maxRight: int) -> int:
            if index == n + 1:
                return 0
            # 花钱
            res = prices[index - 1] + dfs(index + 1, min2(n + 1, 2 * index))
            # 不花钱
            if maxRight >= index:
                res = min2(res, dfs(index + 1, maxRight))
            return res

        n = len(prices)
        res = dfs(1, 0)
        dfs.cache_clear()
        return res


# prices = [1,10,1,1]
