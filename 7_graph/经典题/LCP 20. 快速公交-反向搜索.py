from functools import lru_cache
from typing import List

MOD = int(1e9 + 7)


class Solution:
    def busRapidTransit(
        self, target: int, inc: int, dec: int, jump: List[int], cost: List[int]
    ) -> int:
        """0到target的最少花费  用终点向起点递归
        最终目标 target 的范围在 [1, 10^9]，
        所以仅通过正向 BFS 最坏情况下要遍历全部位置，会超时。所以我们考虑`反向 BFS`
        注意不能遍历全部位置 因此要先坐车再走路

        1 <= jump.length, cost.length <= 10
        1 <= target <= 10^9
        2 <= jump[i] <= 10^6
        1 <= inc, dec, cost[i] <= 10^6
        """

        @lru_cache(None)
        def dfs(pos: int) -> int:
            if pos == 0:
                return 0
            if pos == 1:
                return inc
            res = pos * inc
            for J, C in zip(jump, cost):
                div, mod = divmod(pos, J)
                res = min(res, (dfs(div) + C + inc * mod))
                if mod != 0:
                    res = min(res, dfs(div + 1) + C + dec * (J - mod))

            # 注意这里不能模mod 因为是比大小
            return res

        return dfs(target) % MOD
