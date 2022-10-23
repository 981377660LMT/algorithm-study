# 1000. 合并石头的最低成本
from typing import List
from functools import lru_cache
from itertools import accumulate

# !每次移动（move）需要将连续的 K 堆石头合并为一堆，而这个移动的成本为这 K 堆石头的总数。
# 找出把所有石头合并成一堆的最低成本。如果不可能，返回 -1 。
# 1 <= stones.length <= 30
# 2 <= K <= 30
# 1 <= stones[i] <= 100


INF = int(1e20)


class Solution:
    def mergeStones(self, stones: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            """[left,right] 合并一堆的最低成本"""
            if right - left + 1 < k:
                return 0

            res = INF
            for i in range(left, right, k - 1):  # !左边需要保证(i-left) % (k-1) == 0
                # !最终需要合并成一堆
                mergeCost = preSum[right + 1] - preSum[left] if (right - left) % (k - 1) == 0 else 0
                cand = dfs(left, i) + dfs(i + 1, right) + mergeCost
                if cand < res:
                    res = cand
            return res

        if (len(stones) - k) % (k - 1) != 0:
            return -1
        preSum = list(accumulate(stones, initial=0))
        res = dfs(0, len(stones) - 1)
        dfs.cache_clear()
        return res if res != INF else -1


print(Solution().mergeStones(stones=[3, 2, 4, 1], k=2))
# 输出：20
# 解释：
# 从 [3, 2, 4, 1] 开始。
# 合并 [3, 2]，成本为 5，剩下 [5, 4, 1]。
# 合并 [4, 1]，成本为 5，剩下 [5, 5]。
# 合并 [5, 5]，成本为 10，剩下 [10]。
# 总成本 20，这是可能的最小值。
