from typing import List, Tuple
from heapq import heappop, heappush

# 两个不重叠区间的最小长度和


class Solution:
    def solve(self, events: List[List[int]]) -> int:
        events.sort()
        pq: List[Tuple[int, int]] = []
        res, preMin = int(1e99), int(1e99)

        for start, end in events:
            heappush(pq, (end, end - start + 1))
            while pq and pq[0][0] < start:
                _, curLen = heappop(pq)
                preMin = min(preMin, curLen)
            res = min(res, preMin + end - start + 1)

        return res if res < int(1e20) else 0


print(Solution().solve(events=[[1, 4], [8, 9], [3, 5]]))
# 如果要参加3个，那么pre_max就需要2个数记录
# 如果没有参加限制，则转化为出租车问题dp+二分 参见 1235. 规划兼职工作
