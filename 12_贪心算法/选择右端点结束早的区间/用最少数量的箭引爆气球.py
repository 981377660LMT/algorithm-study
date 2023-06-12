# https://leetcode.cn/problems/minimum-number-of-arrows-to-burst-balloons/
# 452. 用最少数量的箭引爆气球
# 若有一个气球的直径的开始和结束坐标为 xstart，xend，
# !且满足  xstart ≤ x ≤ xend，则该气球会被 引爆
# !返回引爆所有气球所必须射出的 最小 弓箭数 。


# !优先选择结束时间早的区间(覆盖区间的最少点数->区间选点问题)


from typing import List

INF = int(1e18)


def findMinArrowShots(points: List[List[int]]) -> int:
    """在数轴上选择一些点，使得所有区间都包含至少一个点.求最少选择多少个点."""
    points.sort(key=lambda x: x[1])
    preEnd = -INF
    res = 0
    for curStart, curEnd in points:
        if curStart > preEnd:
            res += 1
            preEnd = curEnd
    return res


assert findMinArrowShots([[10, 16], [2, 8], [1, 6], [7, 12]]) == 2
