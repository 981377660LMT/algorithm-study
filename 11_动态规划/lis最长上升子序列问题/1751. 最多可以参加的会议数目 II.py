from typing import List
from functools import lru_cache
from bisect import bisect_right


class Solution:
    def maxValue(self, events: List[List[int]], k: int) -> int:
        # 考虑index之后的会议，选k个的最大value
        @lru_cache(None)
        def dp(index: int, k: int) -> int:
            if k == 0 or index >= len(events):
                return 0
            next = bisect_right(starts, events[index][1])
            print(next, index, k, events)
            # 参加或不参加
            # 如果当前参加了 下一个的start就必须严格大于当前的end 所以使用bisectRight寻找
            return max(events[index][2] + dp(next, k - 1), dp(index + 1, k))

        events.sort()
        starts = [e[0] for e in events]

        return dp(0, k)


print(Solution().maxValue([[1, 2, 4], [3, 4, 3], [2, 3, 1]], 2))
