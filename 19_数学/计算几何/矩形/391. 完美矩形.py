# 391. 完美矩形
# 判断它们是否能精确地覆盖一个矩形区域。
# 每个矩形用左下角的点和右上角的点的坐标来表示。例如， 一个单位正方形可以表示为 [1,1,2,2]

from typing import List

INF = int(1e18)


class Solution:
    def isRectangleCover(self, rectangles: List[List[int]]) -> bool:
        points = set()
        area = 0
        for x1, y1, x2, y2 in rectangles:
            area += (x2 - x1) * (y2 - y1)
            # 合并端点
            points.symmetric_difference_update([(x1, y1), (x1, y2), (x2, y1), (x2, y2)])

        # 最后要矩形必须是四个点
        if len(points) != 4:
            return False

        x1 = INF
        x2 = -INF
        y1 = INF
        y2 = -INF
        for x, y in points:
            x1 = min(x1, x)
            x2 = max(x2, x)
            y1 = min(y1, y)
            y2 = max(y2, y)
        return (x2 - x1) * (y2 - y1) == area
