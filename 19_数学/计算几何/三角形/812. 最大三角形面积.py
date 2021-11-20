# 3 <= points.length <= 50.
# -50 <= points[i][j] <= 50.
# 叉乘的模除表示平行四边形面积

from typing import List
from itertools import combinations


class Solution:
    def largestTriangleArea(self, points: List[List[int]]) -> float:
        def cal_cross_product(A, B, C):
            AB = [B[0] - A[0], B[1] - A[1]]
            AC = [C[0] - A[0], C[1] - A[1]]
            return AB[0] * AC[1] - AB[1] * AC[0]

        return max(abs(cal_cross_product(A, B, C)) / 2 for A, B, C in combinations(points, 3))


print(Solution().largestTriangleArea([[0, 0], [0, 1], [1, 0], [0, 2], [2, 0]]))
