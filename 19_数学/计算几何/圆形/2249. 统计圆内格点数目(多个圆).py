# 2249. 统计圆内格点数目(多个圆)
# https://leetcode.cn/problems/count-lattice-points-inside-a-circle/
# !len(circles) <= 200 1<=r<=x,y<=100

# 注意幂运算不要用** 要用 a*a
# 用pow 和 **都会超时

from collections import defaultdict
from itertools import product
from typing import List


INF = int(1e20)


class Solution:
    def countLatticePoints1(self, circles: List[List[int]]) -> int:
        """O(n^3)遍历点的思路

        优化：
        1.确定点的范围
        2.循环前对circles排序
        """
        res = 0
        xMin, xMax, yMin, yMax = INF, -INF, INF, -INF
        for x, y, r in circles:
            xMin = min(xMin, x - r)
            xMax = max(xMax, x + r)
            yMin = min(yMin, y - r)
            yMax = max(yMax, y + r)

        # 按半径从大到小排序，这样能更早遇到包含 (x,y) 的圆
        circles.sort(key=lambda x: -x[2])
        for x in range(xMin, xMax + 1):
            for y in range(yMin, yMax + 1):
                for cx, cy, r in circles:
                    # 注意要用* 不用**  `** 是logn的操作`
                    if (x - cx) * (x - cx) + (y - cy) * (y - cy) <= r * r:
                        res += 1
                        break
        return res

    def countLatticePoints2(self, circles: List[List[int]]) -> int:
        """O(n^3)遍历圆的思路 怕超时的话可以只计算四分之一个圆内的点，然后映射"""
        res = set()
        for x, y, r in circles:
            for dx, dy in product(range(r + 1), repeat=2):
                if dx * dx + dy * dy <= r * r:
                    res |= {(x + dx, y + dy), (x + dx, y - dy), (x - dx, y + dy), (x - dx, y - dy)}
        return len(res)

    def countLatticePoints3(self, circles: List[List[int]]) -> int:
        """O(n*r)差分的思路：对每个圆 固定y扫描x看那些范围被覆盖
        还原差分数组 对每一个y 扫描  看点是否被覆盖(值>=1)
        """

        diffGroup = defaultdict(lambda: defaultdict(int))
        for cx, cy, cr in circles:
            # 求左右端点，双指针的思想扩展
            dx = 0
            for dy in range(cr, -1, -1):
                while dy * dy + dx * dx <= cr * cr:
                    dx += 1
                # 上半
                diffGroup[cy + dy][cx - dx + 1] += 1
                diffGroup[cy + dy][cx + dx] -= 1
                # 下半对称
                diffGroup[cy - dy][cx - dx + 1] += 1
                diffGroup[cy - dy][cx + dx] -= 1

        res = 0
        for diff in diffGroup.values():
            pre, preSum = -1, 0  # 被几个圆覆盖
            for key in sorted(diff):
                if preSum > 0:
                    res += key - pre
                pre, preSum = key, preSum + diff[key]
        return res
