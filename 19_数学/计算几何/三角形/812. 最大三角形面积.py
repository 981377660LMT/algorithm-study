# 3 <= points.length <= 50.
# -50 <= points[i][j] <= 50.
# 叉乘的模除表示平行四边形面积

from typing import List
from itertools import combinations

Point = List[int]

# https://leetcode.cn/problems/largest-triangle-area/comments/
class Solution:
    def largestTriangleArea(self, points: List[Point]) -> float:
        """叉乘计算三角形面积"""

        def calCross(pa: Point, pb: Point, pc: Point) -> int:
            ab = [pb[0] - pa[0], pb[1] - pa[1]]
            ac = [pc[0] - pa[0], pc[1] - pa[1]]
            return ab[0] * ac[1] - ab[1] * ac[0]

        return max(abs(calCross(pa, pb, pc)) / 2 for pa, pb, pc in combinations(points, 3))


print(Solution().largestTriangleArea([[0, 0], [0, 1], [1, 0], [0, 2], [2, 0]]))
