# 2054. 两个最好的不重叠活动

from bisect import bisect_right
from typing import List, Tuple
from heapq import heappop, heappush

# 你 最多 可以参加 两个时间不重叠 活动，使得它们的价值之和 最大 。
# 请你返回价值之和的 最大值 。
# !如果你参加一个活动，且结束时间为 t ，那么下一个活动必须在 t + 1 或之后的时间开始。
# 你 最多 可以参加 `两个`时间不重叠 活动，使得它们的价值之和 最大 。


# !关键思路：用堆维护event的结束时间,用一个变量维护之前的最大值
# 强化版:11_动态规划/出租车问题/1751. 最多可以参加的会议数目 II.py


def max(x, y):
    if x > y:
        return x
    return y


# events[i] = [startDayi, endDayi, valuei] 。
def maxTwoEvents(events: List[Tuple[int, int, int]]) -> int:
    events = sorted(events, key=lambda x: x[0])
    pq = []  # (end, val)
    res, pre_max = 0, 0
    for start, end, val in events:
        heappush(pq, (end, val))
        while pq and pq[0][0] < start:
            _, pre_val = heappop(pq)
            pre_max = max(pre_max, pre_val)
        res = max(res, pre_max + val)
    return res


# O(nk+nlogn) K个最好的不重叠活动
def maxKEvents(events: List[Tuple[int, int, int]], k: int) -> int:
    n = len(events)
    events = sorted(events, key=lambda x: x[1])
    dp = [[0] * (k + 1) for _ in range(n + 1)]
    for i in range(n):
        start, _, score = events[i]
        dp[i + 1] = dp[i][:]
        prePos = bisect_right(events, start - 1, key=lambda x: x[1]) - 1
        for j in range(1, k + 1):
            dp[i + 1][j] = max(dp[i + 1][j], score + dp[prePos + 1][j - 1])
    return dp[n][k]


print(maxTwoEvents([[1, 3, 2], [4, 5, 2], [2, 4, 3]]))
# 如果要参加3个，那么pre_max就需要2个数记录
# 如果没有参加限制，则转化为出租车问题dp+二分 参见 1235. 规划兼职工作
# 需要按照结束顺序来dp
