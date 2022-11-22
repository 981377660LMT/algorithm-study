from typing import List
from functools import lru_cache

# 每回合，玩家从行的 开始 或 结束 处取走整堆石头。
# 这种情况一直持续到没有更多的石子堆为止，此时手中 石子最多 的玩家 获胜 。


INF = int(1e18)


class Solution:
    def stoneGame(self, piles: List[int]) -> bool:
        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left == right:
                return 0

            res = piles[left] - dfs(left + 1, right)
            res = max(res, piles[right] - dfs(left, right - 1))
            return res

        return dfs(0, len(piles) - 1) > 0


print(Solution().stoneGame(piles=[5, 3, 4, 5]))
