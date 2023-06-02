# Interval Scheduling (weighted)
# 安排会议问题


# 不重叠区间的最大和(出租车问题)
# 给定 N 个闭区间 [ai,bi,scorei]，请你在数轴上选择若干区间，
# !使得选中的区间之间互不相交（端点不可重合）。
# 输出可选取区间的最大权值和。


from bisect import bisect_right
from typing import List, Tuple


def weightedIntervalScheduling(intervals: List[Tuple[int, int, int]]) -> int:
    """选中的区间之间互不相交（端点不可重合）。
    返回可选取区间的最大权值和.
    """
    n = len(intervals)
    intervals = sorted(intervals, key=lambda x: x[1])
    pre = [0] * n
    for i in range(n):
        start = intervals[i][0]
        # 端点不可重合
        pre[i] = bisect_right(intervals, start, key=lambda x: x[1])
    dp = [0] * (n + 1)
    for i in range(n):
        dp[i + 1] = max(dp[i], dp[pre[i]] + intervals[i][2])
    return dp[-1]


if __name__ == "__main__":
    # 435. 无重叠区间
    # 给定一个区间的集合 intervals ，其中 intervals[i] = [starti, endi] 。
    # 返回 需要移除区间的最小数量，使剩余区间互不重叠 。
    class Solution:
        def eraseOverlapIntervals(self, intervals: List[List[int]]) -> int:
            items = [(x[0], x[1], 1) for x in intervals]
            return len(intervals) - weightedIntervalScheduling(items)
