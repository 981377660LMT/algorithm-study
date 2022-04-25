# 注意幂运算不要用** 要用 a*a
# 用pow 和 **都会超时

from itertools import product
from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)

# 优化：
# 1.确定点的范围
# 2.循环前对circles排序


class Solution:
    def countLatticePoints(self, circles: List[List[int]]) -> int:
        """O(n^3)遍历点的思路"""
        res = 0
        xMin, xMax, yMin, yMax = 210, 0, 210, 0
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
        """O(n^2)差分的思路：对每个圆 固定y扫描x看那些范围被覆盖
        还原差分数组 对每一个y 扫描  看点是否被覆盖(值>=1)
        """
        diff = [[0] * 205 for _ in range(205)]
        for cx, cy, cr in circles:
            # 求左右端点，双指针的思想扩展
            dx = 0
            for dy in range(cr, -1, -1):
                while dy * dy + dx * dx <= cr * cr:
                    dx += 1
                # 上半
                diff[cy + dy][cx - dx + 1] += 1
                diff[cy + dy][cx + dx] -= 1
                # 下半对称
                diff[cy - dy][cx - dx + 1] += 1
                diff[cy - dy][cx + dx] -= 1

        res = 0
        for y in range(205):
            curSum = 0
            for x in range(205):
                curSum += diff[y][x]
                if curSum >= 1:
                    res += 1
        return res
