# !最大不相交区间数量(maxNonIntersectingIntervals)
# 给定 N 个闭区间 [ai,bi]，请你在数轴上选择若干区间，
# 使得选中的区间之间互不相交（包括端点）。
# 输出可选取区间的最大数量，返回选择方案.

# 1. 右端点从小到大排序(选择结束时间早的(罗志祥贪心问题))
# 2. 遍历区间，如果已经包含点，pass，否则`选择当前区间右端点`
################################################


from typing import List, Tuple, Union


def maxNonIntersectingIntervals(
    intervals: Union[List[Tuple[int, int]], List[List[int]]],
    allowOverlapping=False,
    endInclusive=True,
) -> List[int]:
    """
    给定 n 个区间 [left_i,right_i].
    请你在数轴上选择若干区间,使得选中的区间之间互不相交.

    Args:
        intervals: 区间列表,每个区间为[left,right].
        allowOverlapping: 是否允许选择的区间端点重合.默认为False.
        endInclusive: 是否包含区间右端点.默认为True.

    Returns:
        List[int]: 区间索引列表.如果有多种方案，返回字典序最小的那个.
    """

    n = len(intervals)
    if n == 0:
        return []
    if n == 1:
        return [0]
    order = sorted(range(n), key=lambda x: intervals[x][1])
    res = [order[0]]
    preEnd = intervals[order[0]][1]
    for i in order[1:]:
        start, end = intervals[i]
        if not endInclusive:
            end -= 1
        if allowOverlapping:
            if start >= preEnd:
                res.append(i)
                preEnd = end
        else:
            if start > preEnd:
                res.append(i)
                preEnd = end
    return res


if __name__ == "__main__":

    class Solution:
        # 435. 无重叠区间
        # https://leetcode.cn/problems/non-overlapping-intervals/description/
        # 给定一个区间的集合 intervals ，其中 intervals[i] = [starti, endi] 。
        # 返回 需要移除区间的最小数量，使剩余区间互不重叠 。
        def eraseOverlapIntervals(self, intervals: List[List[int]]) -> int:
            res = maxNonIntersectingIntervals(intervals, allowOverlapping=True, endInclusive=True)
            return len(intervals) - len(res)

        # 452. 用最少数量的箭引爆气球
        # https://leetcode-cn.com/problems/minimum-number-of-arrows-to-burst-balloons/
        # 给定气球在一个水平数轴上的坐标。如果有n个气球，它们的坐标分别是x1，x2，...，xn。
        # 如果一支箭能引爆气球，则意味着这支箭的位置x，存在x在x1和x2之间，或者x在x2和x3之间，或者x在x3和x4之间，依此类推。
        # 任何有重叠的气球都可以通过一支箭被引爆。
        # 返回所需的最小射箭数量。
        def findMinArrowShots(self, points: List[List[int]]) -> int:
            res = maxNonIntersectingIntervals(points, allowOverlapping=False, endInclusive=True)
            return len(res)

        # 646. 最长数对链
        # https://leetcode.cn/problems/maximum-length-of-pair-chain/
        def findLongestChain(self, pairs: List[List[int]]) -> int:
            res = maxNonIntersectingIntervals(pairs, allowOverlapping=False, endInclusive=True)
            return len(res)

        # 100657. 不相交子字符串的最大数量
        # https://leetcode.cn/contest/biweekly-contest-157/problems/find-maximum-number-of-non-intersecting-substrings/
        # !返回以 首尾字母相同 且 长度至少为 4 的 不相交子字符串 的最大数量。
        def maxSubstrings(self, word: str, k=4) -> int:
            from collections import defaultdict
            from bisect import bisect_left

            mp = defaultdict(list)
            for i, v in enumerate(word):
                mp[v].append(i)

            intervals = []
            for group in mp.values():
                for start in range(len(group)):
                    end = bisect_left(group, group[start] + k - 1, lo=start + 1)
                    if end < len(group):
                        intervals.append((group[start], group[end]))

            return len(maxNonIntersectingIntervals(intervals))
