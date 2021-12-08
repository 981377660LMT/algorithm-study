from typing import List
from functools import lru_cache

# 2 <= locations.length <= 100
# locations[i] 表示第 i 个城市的位置
# 给你 start，finish 和 fuel 分别表示出发城市、目的地城市和你初始拥有的汽油总量

# 每一步中，如果你在城市 i ，你可以选择任意一个城市 j ，满足  j != i 且 0 <= j < locations.length ，并移动到城市 j
# 从城市 i 移动到 j 消耗的汽油量为 |locations[i] - locations[j]|

# 请注意， fuel 任何时刻都 不能 为负，且你 `可以 经过任意城市超过一次`（包括 start 和 finish ）。
# 请你返回从 start 到 finish 所有可能路径的数目。

MOD = int(1e9 + 7)


class Solution:
    def countRoutes(self, locations: List[int], start: int, finish: int, fuel: int) -> int:
        weight = lambda i, j: abs(locations[i] - locations[j])

        @lru_cache(None)
        def dfs(cur: int, remain: int) -> int:
            # if remain < 0:
            #     return 0
            if remain < weight(cur, finish):
                return 0

            res = 0
            if cur == finish:
                res += 1

            for next in range(len(locations)):
                if cur == next:
                    continue
                res += dfs(next, remain - weight(cur, next))

            return res

        return dfs(start, fuel) % MOD


print(Solution().countRoutes(locations=[2, 3, 6, 8, 4], start=1, finish=3, fuel=5))
# 输出：4
# 解释：以下为所有可能路径，每一条都用了 5 单位的汽油：
# 1 -> 3
# 1 -> 2 -> 3
# 1 -> 4 -> 3
# 1 -> 4 -> 2 -> 3
print(Solution().countRoutes(locations=[4, 3, 1], start=1, finish=0, fuel=6))
