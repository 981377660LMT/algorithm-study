from typing import List
from functools import lru_cache
from bisect import bisect_right

# 最多参加k个会议
# 你必须完整参加会议 请你返回能得到的会议价值 最大和 。
# !你不能同时参加一个开始日期与另一个结束日期相同的两个会议
# 1 <= k * events.length <= 1e6
# 1 <= startDayi <= endDayi <= 1e9


class Solution:
    def maxValue(self, events: List[List[int]], k: int) -> int:
        """记忆化dfs"""

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain == 0 or index >= n:
                return 0

            res = dfs(index + 1, remain)  # 不选当前会议
            _, end, score = events[index]  # 选当前会议
            nextPos = bisect_right(events, end, key=lambda x: x[0])
            nextRes = dfs(nextPos, remain - 1) + score
            return res if res > nextRes else nextRes

        n = len(events)
        events.sort(key=lambda x: x[0])  # 按照start排序
        res = dfs(0, k)
        dfs.cache_clear()
        return res

    def maxValue2(self, events: List[List[int]], k: int) -> int:
        """dp"""
        n = len(events)
        events.sort(key=lambda x: x[1])  # 按照end排序
        dp = [[0] * (k + 1) for _ in range(n + 1)]
        for i in range(n):
            start, _, score = events[i]
            dp[i + 1] = dp[i][:]  # 不选当前会议
            # !选当前会议,由于端点不可重合,所以要start-1
            prePos = bisect_right(events, start - 1, key=lambda x: x[1]) - 1
            for j in range(1, k + 1):
                dp[i + 1][j] = max(dp[i + 1][j], score + dp[prePos + 1][j - 1])

        return dp[-1][-1]


print(Solution().maxValue([[1, 2, 4], [3, 4, 3], [2, 3, 1]], 2))
print(Solution().maxValue2([[1, 2, 4], [3, 4, 3], [2, 3, 1]], 2))
print(Solution().maxValue([[1, 2, 4], [3, 4, 3], [2, 3, 10]], 2))
print(Solution().maxValue2([[1, 2, 4], [3, 4, 3], [2, 3, 10]], 2))
