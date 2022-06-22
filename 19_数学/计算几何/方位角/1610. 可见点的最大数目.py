# 计算角度

from typing import List
from math import atan2, degrees
from bisect import bisect_right

# 输出：3
# 解释：阴影区域代表你的视野。在你的视野中，所有的点都清晰可见，尽管 [2,2] 和 [3,3]在同一条直线上，你仍然可以看到 [3,3] 。

# 你的视野就是角度范围 [d - angle/2, d + angle/2] 所指示的那片区域。
# 1 <= points.length <= 105
# 0 <= angle < 360


class Solution:
    def visiblePoints(self, points: List[List[int]], angle: int, location: List[int]) -> int:
        same = 0
        angles = []
        for p in points:
            if p == location:
                same += 1
            else:
                cur = degrees(atan2(p[1] - location[1], p[0] - location[0]))
                angles.extend([cur, cur + 360])  # !注意这个360 转一圈

        angles.sort()
        res = 0
        for i, a in enumerate(angles):
            pos = bisect_right(angles, a + angle) - 1
            if pos - i + 1 > res:
                res = pos - i + 1

        return res + same


# print(Solution().visiblePoints(points=[[2, 1], [2, 2], [3, 3]], angle=90, location=[1, 1]))
# print(Solution().visiblePoints(points=[[1, 0], [2, 1]], angle=13, location=[1, 1]))
print(Solution().visiblePoints(points=[[0, 0], [0, 2]], angle=90, location=[1, 1]))

# atan2(y,x) 函数可以自动判断象限 :方位角,返回以弧度表示的 y/x 的反正切
# atan2函数返回的是原点至点(x,y)的方位角，即与 x 轴的夹角。
# 也可以理解为复数 x+yi 的辐角。返回值的单位为弧度，取值范围为(-pi,pi]
print('test'.center(20, '-'))
print(atan2(1, 1))
print(atan2(1, 0))
print(degrees(atan2(1, 0)))
print(atan2(-1, 1))
print(atan2(0, 0))

