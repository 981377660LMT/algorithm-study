# 给你一个以 (radius, x_center, y_center) 表示的圆和一个与坐标轴平行的矩形 (x1, y1, x2, y2)，
# 其中 (x1, y1) 是矩形左下角的坐标，(x2, y2) 是右上角的坐标。

# 找到矩形中距离圆心最近的位置，研究这个位置是否超过了圆的半径
class Solution:
    def checkOverlap(
        self, radius: int, x_center: int, y_center: int, x1: int, y1: int, x2: int, y2: int
    ) -> bool:
        if x1 <= x_center <= x2 and y1 <= y_center <= y2:
            return True

        # 矩形边上，距离圆心最近的点
        yp = y1 if y_center < y1 else min(y2, y_center)
        xp = x1 if x_center < x1 else min(x2, x_center)

        if (xp - x_center) * (xp - x_center) + (yp - y_center) * (yp - y_center) <= radius * radius:
            return True
        return False


print(Solution().checkOverlap(radius=1, x_center=0, y_center=0, x1=1, y1=-1, x2=3, y2=1))
# 输出：true
# 解释：圆和矩形有公共点 (1,0)

