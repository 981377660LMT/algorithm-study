from typing import List
from functools import lru_cache


class Solution:
    def canCross(self, stones: List[int]) -> bool:
        end = stones[-1]
        available = set(stones)

        @lru_cache(None)
        def dfs(index: int, pre_step: int) -> bool:
            if index == end:
                return True
            for next_step in [pre_step - 1, pre_step, pre_step + 1]:
                if next_step <= 0:
                    continue
                if next_step + index in available and dfs(next_step + index, next_step):
                    return True
            return False

        return dfs(0, 0)

