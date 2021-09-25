from typing import List

#  判断它们是否能精确地覆盖一个矩形区域。
# 每个矩形用左下角的点和右上角的点的坐标来表示。例如， 一个单位正方形可以表示为 [1,1,2,2]
class Solution:
    def isRectangleCover(self, rectangles: List[List[int]]) -> bool:
        points = set()
        area = 0
        for [x1, y1, x2, y2] in rectangles:
            area += (x2 - x1) * (y2 - y1)
            # 合并端点
            points.symmetric_difference_update([(x1, y1), (x1, y2), (x2, y1), (x2, y2)])

        # 最后要矩形必须是四个点
        if len(points) != 4:
            return False
        x1 = float('inf')
        x2 = float('-inf')
        y1 = float('inf')
        y2 = float('-inf')
        for [x, y] in points:
            x1 = min(x1, x)
            x2 = max(x2, x)
            y1 = min(y1, y)
            y2 = max(y2, y)
        return (x2 - x1) * (y2 - y1) == area


# rectangles = [
#   [1,1,3,3],
#   [3,1,4,2],
#   [3,2,4,4],
#   [1,3,2,4],
#   [2,3,3,4]
# ]

# 从「面积」和「顶点」两个维度来判断：
# 1、判断面积，通过完美矩形的理论坐标计算出一个理论面积，然后和 rectangles 中小矩形的实际面积和做对比。
# 2、判断顶点，points 集合中应该只剩下 4 个顶点且剩下的顶点必须都是完美矩形的理论顶点。
