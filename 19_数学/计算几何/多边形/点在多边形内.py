# 判断点是否在多边形上/内
class Solution:
    def solve(self, polygon, x, y):
        assert len(polygon) >= 2, "not a polygon"
