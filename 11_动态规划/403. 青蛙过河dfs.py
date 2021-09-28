from typing import List
from functools import lru_cache


class Solution:
    def canCross(self, stones: List[int]) -> bool:
        end = stones[-1]
        s = set(stones)

        @lru_cache(None)
        def dfs(stone: int, step: int) -> bool:
            if stone == end:
                return True
            for next_step in [step - 1, step, step + 1]:
                if next_step <= 0:
                    continue
                if next_step + stone in s and dfs(next_step + stone, next_step):
                    return True
            return False

        return dfs(0, 0)

