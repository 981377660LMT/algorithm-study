# 1401. 圆和矩形是否有重叠
# https://leetcode.cn/problems/circle-and-rectangle-overlapping/
# 给你一个以 (radius, x_center, y_center) 表示的圆和一个与坐标轴平行的矩形 (x1, y1, x2, y2)，
# 其中 (x1, y1) 是矩形左下角的坐标，(x2, y2) 是右上角的坐标。
#
# !找到矩形中距离圆心最近的位置，研究这个位置是否超过了圆的半径


class Solution:
    def checkOverlap(
        self, radius: int, xCenter: int, yCenter: int, x1: int, y1: int, x2: int, y2: int
    ) -> bool:
        # 圆心到矩形在 x 方向上的最近距离 dx
        if xCenter < x1:
            dx = x1 - xCenter
        elif xCenter > x2:
            dx = xCenter - x2
        else:
            dx = 0

        # y 方向上的最近距离 dy
        if yCenter < y1:
            dy = y1 - yCenter
        elif yCenter > y2:
            dy = yCenter - y2
        else:
            dy = 0

        return dx * dx + dy * dy <= radius * radius


print(Solution().checkOverlap(radius=1, xCenter=0, yCenter=0, x1=1, y1=-1, x2=3, y2=1))
# 输出：true
# 解释：圆和矩形有公共点 (1,0)
