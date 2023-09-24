from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 假设你是一家合金制造公司的老板，你的公司使用多种金属来制造合金。现在共有 n 种不同类型的金属可以使用，并且你可以使用 k 台机器来制造合金。每台机器都需要特定数量的每种金属来创建合金。

# 对于第 i 台机器而言，创建合金需要 composition[i][j] 份 j 类型金属。最初，你拥有 stock[i] 份 i 类型金属，而每购入一份 i 类型金属需要花费 cost[i] 的金钱。

# 给你整数 n、k、budget，下标从 1 开始的二维数组 composition，两个下标从 1 开始的数组 stock 和 cost，请你在预算不超过 budget 金钱的前提下，最大化 公司制造合金的数量。

# 所有合金都需要由同一台机器制造。


# 返回公司可以制造的最大合金数。


class Solution:
    def maxNumberOfAlloys(
        self,
        n: int,
        k: int,
        budget: int,
        composition: List[List[int]],
        stock: List[int],
        cost: List[int],
    ) -> int:
        res = 0

        # 枚举机器
        for i in range(k):
            need = composition[i]

            # 二分最大合金
            def check(mid: int) -> bool:
                """造出 mid 个合金需要多少金钱"""
                allCost = 0
                for j in range(n):
                    allCost += max(0, need[j] * mid - stock[j]) * cost[j]
                return allCost <= budget

            left, right = 1, int(1e9)
            ok = False
            while left <= right:
                mid = (left + right) // 2
                if check(mid):
                    left = mid + 1
                    ok = True
                else:
                    right = mid - 1
            res = max(res, right)

        return res
