# 爱丽丝和鲍勃继续他们的石子游戏。许多堆石子 排成一行，每堆都有正整数颗石子 piles[i]。
# 游戏以谁手中的石子最多来决出胜负。
# 爱丽丝和鲍勃轮流进行，爱丽丝先开始。最初，M = 1。
# !在每个玩家的回合中，该玩家可以拿走剩下的 前 X 堆的所有石子，其中 1 <= X <= 2M。然后，令 M = max(M, X)。
# 游戏一直持续到所有石子都被拿走。
# 假设爱丽丝和鲍勃都发挥出最佳水平，返回爱丽丝可以得到的最大数量的石头。


from typing import List
from functools import lru_cache
from itertools import accumulate

INF = int(1e20)


class Solution:
    def stoneGameII(self, piles: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, M: int) -> int:
            # 可以一口气拿完剩下的石头
            if index + 2 * M >= len(sufSum):
                return sufSum[index]

            res = -INF
            for i in range(1, 2 * M + 1):
                # 减去对手最多拿的
                res = max(res, sufSum[index] - dfs(index + i, max(i, M)))
            return res

        sufSum = list(accumulate(piles[::-1]))[::-1]
        return dfs(0, 1)


print(Solution().stoneGameII(piles=[2, 7, 9, 4, 4]))
