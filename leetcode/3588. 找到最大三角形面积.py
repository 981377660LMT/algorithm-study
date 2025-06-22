# 3588. 找到最大三角形面积
# https://leetcode.cn/problems/find-maximum-area-of-a-triangle/
#
# 给你一个二维数组 coords，大小为 n x 2，表示一个无限笛卡尔平面上 n 个点的坐标。
# 找出一个 最大 三角形的 两倍 面积，其中三角形的三个顶点来自 coords 中的任意三个点，并且该三角形至少有一条边与 x 轴或 y 轴平行。
# 严格地说，如果该三角形的最大面积为 A，则返回 2 * A。
# 如果不存在这样的三角形，返回 -1。
# 注意，三角形的面积 不能 为零。
#
# !1.讨论三角形底边与 y 轴平行的情况，x 轴平行的情况可以通过交换坐标来处理。


from collections import defaultdict
from typing import List

INF = int(1e18)


class Solution:
    def maxArea(self, coords: List[List[int]]) -> int:
        def calc(points: List[List[int]]) -> int:
            """按照x分组, 维护最小/最大y值."""
            minX, maxX = INF, -INF
            xToMinY = defaultdict(lambda: INF)
            xToMaxY = defaultdict(lambda: -INF)
            for x, y in points:
                minX = min(minX, x)
                maxX = max(maxX, x)
                xToMinY[x] = min(xToMinY[x], y)
                xToMaxY[x] = max(xToMaxY[x], y)

            res = 0
            for x, y in xToMinY.items():
                tmp = max(x - minX, maxX - x)
                res = max(res, tmp * (xToMaxY[x] - y))
            return res

        res1 = calc(coords)
        res2 = calc([[y, x] for x, y in coords])
        res = max(res1, res2)
        return res if res > 0 else -1
