from typing import List
from math import atan2, pi
from bisect import bisect_right

# 输出：3
# 解释：阴影区域代表你的视野。在你的视野中，所有的点都清晰可见，尽管 [2,2] 和 [3,3]在同一条直线上，你仍然可以看到 [3,3] 。

# 你的视野就是角度范围 [d - angle/2, d + angle/2] 所指示的那片区域。
# 1 <= points.length <= 105


class Solution:
    def visiblePoints(self, points: List[List[int]], angle: int, location: List[int]) -> int:
        sameCount = 0
        polarDegrees = []
        for p in points:
            if p == location:
                sameCount += 1
            else:
                polarDegrees.append(atan2(p[1] - location[1], p[0] - location[0]))

        polarDegrees.sort()
        # 点是循环的,要转一圈
        polarDegrees += [d + 2 * pi for d in polarDegrees]
        print(polarDegrees)

        n = len(polarDegrees)
        curDegree = angle * pi / 180
        maxCount = max(
            [bisect_right(polarDegrees, polarDegrees[i] + curDegree) - i for i in range(n)],
            default=0,
        )

        return maxCount + sameCount


print(Solution().visiblePoints(points=[[2, 1], [2, 2], [3, 3]], angle=90, location=[1, 1]))

# atan2(y,x) 函数可以自动判断象限 :方位角,返回以弧度表示的 y/x 的反正切
# atan2函数返回的是原点至点(x,y)的方位角，即与 x 轴的夹角。
# 也可以理解为复数 x+yi 的辐角。返回值的单位为弧度，取值范围为(-pi,pi]
print(atan2(1, 1))
print(atan2(-1, 1))
print(atan2(0, 0))

