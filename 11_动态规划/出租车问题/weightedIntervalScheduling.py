# Interval Scheduling (weighted)
# 安排会议问题
# 不重叠区间的最大和(出租车问题)
# 给定 N 个闭区间 [ai,bi,scorei]，请你在数轴上选择若干区间，
# !使得选中的区间之间互不相交（端点不可重合）。
# 输出可选取区间的最大权值和。


from typing import List, Tuple
from bisect import bisect_left, bisect_right


def weightedIntervalScheduling(intervals: List[Tuple[int, int, int]], overlap=False) -> int:
    """
    给定 n 个闭区间 [left_i,right_i,score_i].
    请你在数轴上选择若干区间,使得选中的区间之间互不相交.
    返回可选取区间的最大权值和.

    Args:
        intervals: 区间列表,每个区间为[left,right,score].
        overlapping: 是否允许选择的区间端点重合.默认为False.
    """
    n = len(intervals)
    intervals = sorted(intervals, key=lambda x: x[1])
    pre = [0] * n
    f = bisect_right if overlap else bisect_left
    for i in range(n):
        start = intervals[i][0]
        pre[i] = f(intervals, start, key=lambda x: x[1])
    dp = [0] * (n + 1)
    for i in range(n):
        dp[i + 1] = max(dp[i], dp[pre[i]] + intervals[i][2])
    return dp[n]


if __name__ == "__main__":
    # 435. 无重叠区间
    # 给定一个区间的集合 intervals ，其中 intervals[i] = [starti, endi] 。
    # 返回 需要移除区间的最小数量，使剩余区间互不重叠 。
    class Solution:
        def eraseOverlapIntervals(self, intervals: List[List[int]]) -> int:
            items = [(x[0], x[1], 1) for x in intervals]
            return len(intervals) - weightedIntervalScheduling(items, True)  # type: ignore

    # 2830. 销售利润最大化
    # 给你一个整数 n 表示数轴上的房屋数量，编号从 0 到 n - 1 。
    # 另给你一个二维整数数组 offers ，其中 offers[i] = [starti, endi, goldi] 表示第 i 个买家想要以 goldi 枚金币的价格购买从 starti 到 endi 的所有房屋。
    # 作为一名销售，你需要有策略地选择并销售房屋使自己的收入最大化。
    # 返回你可以赚取的金币的最大数目。
    # 注意 同一所房屋不能卖给不同的买家，并且允许保留一些房屋不进行出售。

    class Solution2:
        def maximizeTheProfit(self, n: int, offers: List[List[int]]) -> int:
            return weightedIntervalScheduling(offers)  # type: ignore

    # 2008. 出租车的最大盈利
    # https://leetcode.cn/problems/maximum-earnings-from-taxi/
    class Solution3:
        def maxTaxiEarnings(self, n: int, rides: List[List[int]]) -> int:
            intervals = [(x[0], x[1], x[1] - x[0] + x[2]) for x in rides]
            return weightedIntervalScheduling(intervals, overlap=True)  # type: ignore

    # 1235. 规划兼职工作
    class Solution4:
        def jobScheduling(self, startTime: List[int], endTime: List[int], profit: List[int]) -> int:
            intervals = list(tuple(zip(startTime, endTime, profit)))
            return weightedIntervalScheduling(intervals, overlap=True)
