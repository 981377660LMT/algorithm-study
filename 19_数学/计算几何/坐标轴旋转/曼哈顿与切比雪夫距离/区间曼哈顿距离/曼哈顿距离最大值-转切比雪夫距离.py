# https://leetcode.cn/problems/minimize-manhattan-distances/description/
# 3102. 最小化曼哈顿距离
# 给你一个下标从 0 开始的数组 points ，它表示二维平面上一些点的整数坐标，其中 points[i] = [xi, yi] 。
# 两点之间的距离定义为它们的曼哈顿距离。
# 请你恰好移除一个点，返回移除后任意两点之间的 最大 距离可能的 最小 值。
#
# 任意两点曼哈顿距离的最大值.
# !我们只需要计算任意两个点变换后 (x,y) -> (y+x,y-x) 切比雪夫距离的最大值，
# !即横纵坐标差的最大值 max(abs(x1-x2),abs(y1-y2)).


from typing import List

from sortedcontainers import SortedList

INF = int(1e18)


class Solution:
    def minimumDistance(self, points: List[List[int]]) -> int:
        xs, ys = SortedList(), SortedList()
        for x, y in points:
            xs.add(y + x)
            ys.add(y - x)
        res = INF
        for x, y in points:
            xx, yy = y + x, y - x
            xs.remove(xx)
            ys.remove(yy)
            res = min(res, max(abs(xs[-1] - xs[0]), abs(ys[-1] - ys[0])))
            xs.add(xx)
            ys.add(yy)
        return res
