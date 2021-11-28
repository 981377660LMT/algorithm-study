# 每次移动（move）需要将连续的 K 堆石头合并为一堆，而这个移动的成本为这 K 堆石头的总数。
# 找出把所有石头合并成一堆的最低成本。如果不可能，返回 -1 。
from typing import List
from functools import lru_cache
from itertools import accumulate

INF = 0x7FFFFFFF

# dfs涉及三个维度的状态 每个状态都要遍历一遍 所以时间复杂度O(n^3)


class Solution:
    def mergeStones(self, stones: List[int], k: int) -> int:
        prefix = [0] + list(accumulate(stones))

        # dp[i][j][p] means the cost needed to merge stone[i] ~ stones[j] into p piles.
        @lru_cache(None)
        def dfs(l: int, r: int, p: int) -> int:
            if ((r - l + 1) - p) % (k - 1) != 0:
                return INF
            if l == r:
                return 0 if p == 1 else INF
            if p == 1:
                return dfs(l, r, k) + prefix[r + 1] - prefix[l]

            return min(dfs(l, mid, 1) + dfs(mid + 1, r, p - 1) for mid in range(l, r, k - 1))

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
