# 判断点是否在多边形上/内
# Ray casting algorithm
class Solution:
    def solve(self, polygon, x, y):
        assert len(polygon) >= 3, "not a polygon"
        n = len(polygon)
        is_inside = False

        for i in range(n):
            x0, y0 = polygon[i]
            x1, y1 = polygon[(i + 1) % n]

            if not min(y0, y1) < y <= max(y0, y1):
                continue

            slope = (x1 - x0) / (y1 - y0)
            x_i = x0 + (y - y0) * slope

            if x_i < x:
                is_inside = not is_inside

        return is_inside


print(Solution().solve(polygon=[[-3, -3], [-3, 3], [3, 3], [3, -3]], x=0, y=0))
# 我们取一条从多边形外部开始，以给定目标坐标为终点的射线，
# 并计算该射线与多边形边之间的交点数。每次光线与边相交时，
# 我们要么进入多边形，要么离开它。
# 因此，奇数交集计数表示我们在多边形内部，偶数表示我们在外部。
