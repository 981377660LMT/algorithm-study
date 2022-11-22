# Alice 和 Bob 交替进行一个游戏，由 Alice 先手。
# 在游戏中，共有 n 堆石头。在每个玩家的回合中，玩家需要 选择 任一非空石头堆，
# 从中移除任意 非零 数量的石头。如果不能移除任意的石头，就输掉游戏，同时另一人获胜。
# 给定一个整数数组 piles ，piles[i] 为 第 i 堆石头的数量，如果 Alice 能获胜返回 true ，反之返回 false 。
# Alice 和 Bob 都会采取 最优策略 。
# 1 <= piles[i] <= 7

# !异或结果不等于 0 时，先手必胜

from typing import List, Tuple
from functools import lru_cache, reduce


class Solution:
    def nimGame2(self, piles: List[int]) -> bool:
        return reduce(lambda x, y: x ^ y, piles, 0) != 0

    def nimGame(self, piles: List[int]) -> bool:
        @lru_cache(None)
        def dfs(nums: Tuple[int]) -> bool:
            if nums == END:
                return False

            counter = list(nums)
            for i, count in enumerate(counter):
                if count == 0:
                    continue
                for remain in range(count):
                    counter[i] = remain
                    if not dfs(tuple(counter)):
                        return True
                    counter[i] = count

            return False

        n = len(piles)
        END = tuple([0] * n)
        return dfs(tuple(piles))


print(Solution().nimGame([1, 2, 3]))
