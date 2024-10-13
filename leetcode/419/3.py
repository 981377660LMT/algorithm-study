from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


# Alice 和 Bob 正在玩一个幻想战斗游戏，游戏共有 n 回合，每回合双方各自都会召唤一个魔法生物：火龙（F）、水蛇（W）或地精（E）。每回合中，双方 同时 召唤魔法生物，并根据以下规则得分：

# 如果一方召唤火龙而另一方召唤地精，召唤 火龙 的玩家将获得一分。
# 如果一方召唤水蛇而另一方召唤火龙，召唤 水蛇 的玩家将获得一分。
# 如果一方召唤地精而另一方召唤水蛇，召唤 地精 的玩家将获得一分。
# 如果双方召唤相同的生物，那么两个玩家都不会获得分数。
# 给你一个字符串 s，包含 n 个字符 'F'、'W' 和 'E'，代表 Alice 每回合召唤的生物序列：

# 如果 s[i] == 'F'，Alice 召唤火龙。
# 如果 s[i] == 'W'，Alice 召唤水蛇。
# 如果 s[i] == 'E'，Alice 召唤地精。

# Bob 的出招序列未知，但保证 Bob 不会在连续两个回合中召唤相同的生物。如果在 n 轮后 Bob 获得的总分 严格大于 Alice 的总分，则 Bob 战胜 Alice。

# 返回 Bob 可以用来战胜 Alice 的不同出招序列的数量。


# 由于答案可能非常大，请返回答案对 109 + 7 取余 后的结果。


mp = {"F": 0, "W": 1, "E": 2}
winBy = [1, 2, 0]


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def countWinningSequences(self, s: str) -> int:
        nums = [mp[v] for v in s]
        winBys = [winBy[v] for v in nums]

        @lru_cache(None)
        def dfs(index: int, preType: int, score: int) -> int:
            if index == n:
                return score > 0
            res = 0
            for i in range(3):
                if i != preType:
                    cost = 0
                    if winBys[index] == i:
                        cost = 1
                    elif winBy[i] == nums[index]:
                        cost = -1
                    res += dfs(index + 1, i, score + cost)
            return res % MOD

        n = len(s)
        res = dfs(0, -1, 0)
        dfs.cache_clear()
        return res % MOD


# s = "FWEFW"

print(Solution().countWinningSequences(s="FWEFW"))
