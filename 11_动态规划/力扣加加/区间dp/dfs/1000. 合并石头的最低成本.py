# 每次移动（move）需要将连续的 K 堆石头合并为一堆，而这个移动的成本为这 K 堆石头的总数。
# 找出把所有石头合并成一堆的最低成本。如果不可能，返回 -1 。
from typing import List
from functools import lru_cache
from itertools import accumulate

INF = 0x7FFFFFFF

# dfs涉及三个维度的状态 每个状态都要遍历一遍 所以时间复杂度O(n^3)
# dp[i][j][m] means the cost needed to merge stone[i] ~ stones[j] into m piles.

# Initial status dp[i][i][1] = 0 and dp[i][i][m] = infinity

# dp[i][j][1] = dp[i][j][k] + stonesNumber[i][j]
# dp[i][j][m] = min(dp[i][mid][1] + dp[mid + 1][j][m - 1])


class Solution:
    def mergeStones(self, stones: List[int], k: int) -> int:
        prefix = [0] + list(accumulate(stones))

        @lru_cache(None)
        def dfs(left: int, right: int, targetPile: int) -> int:
            if ((right - left + 1) - targetPile) % (k - 1) != 0:
                return INF
            if left == right:
                return 0 if targetPile == 1 else INF
            if targetPile == 1:
                return dfs(left, right, k) + prefix[right + 1] - prefix[left]

            res = 0x7FFFFFFF
            for mid in range(left, right, k - 1):
                res = min(res, dfs(left, mid, 1) + dfs(mid + 1, right, targetPile - 1))
            return res

        res = dfs(0, len(stones) - 1, 1)
        return res if res != INF else -1


print(Solution().mergeStones(stones=[3, 2, 4, 1], k=2))
# 输出：20
# 解释：
# 从 [3, 2, 4, 1] 开始。
# 合并 [3, 2]，成本为 5，剩下 [5, 4, 1]。
# 合并 [4, 1]，成本为 5，剩下 [5, 5]。
# 合并 [5, 5]，成本为 10，剩下 [10]。
# 总成本 20，这是可能的最小值。
