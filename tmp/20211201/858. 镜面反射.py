# 正方形房间的墙壁长度为 col，一束激光从西南角射出，首先会与东墙相遇，入射点到接收器 0 的距离为 q 。
# 返回光线最先遇到的接收器的编号（保证光线最终会遇到一个接收器）。

# (光线反射/镜面反射)
# !思路：求出最小公倍数，算相遇时在两个维度走的步数奇偶性
from math import lcm


class Solution:
    def mirrorReflection(self, col: int, row: int) -> int:
        lcm_ = lcm(col, row)
        colCount, rowCount = lcm_ // col, lcm_ // row
        if rowCount & 1 and colCount & 1:
            return 1
        if rowCount & 1:
            return 0
        return 2


print(Solution().mirrorReflection(col=2, row=1))
