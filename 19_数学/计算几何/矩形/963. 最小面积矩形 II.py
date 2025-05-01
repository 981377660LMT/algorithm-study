# 给定在 xy 平面上的一组点，确定由这些点组成的任何矩形的最小面积，其中矩形的边不一定平行于 x 轴和 y 轴。
# 如果没有任何矩形，就返回 0。
#
# 1 <= points.length <= 50
# https://leetcode.cn/problems/minimum-area-rectangle-ii/solutions/707666/c-on2-0ms-100-by-hqztrue-9ij7/
#
# 最坏情况下n个点可以组成Θ(n^2logn)个矩形

from typing import List
from collections import defaultdict

INF = int(1e18)


class Solution:
    def minAreaFreeRect(self, points: List[List[int]]) -> float:
        """
        在 arbitrary orientation 下寻找最小面积矩形。
        核心思路：任何矩形的两条对角线相交且等长。我们枚举所有点对作“对角线”，
        以 (mid_x, mid_y, dist²) 为 key 分组；同组内任意两条对角线对应一个矩形。
        矩形面积＝以某对角线端点为顶点的两条邻边的叉积绝对值。
        时间复杂度 O(n^2·g)（g 是平均组内大小，点数 ≤50 时足够快）。
        """
        n = len(points)
        # key: (sum_x, sum_y, dist2) -> list of pairs of indices (i,j)
        diag_map = defaultdict(list)
        res = INF

        for i in range(n):
            x1, y1 = points[i]
            for j in range(i + 1, n):
                x2, y2 = points[j]
                # 对角线中点 * 2（避免浮点运算），以及对角线长的平方
                key = (x1 + x2, y1 + y2, (x1 - x2) ** 2 + (y1 - y2) ** 2)

                # 若同一 key 下已有先前对角线，任取一条先前的 (k,l)
                # 与当前 (i,j) 构成矩形
                for k, l in diag_map[key]:
                    xk, yk = points[k]
                    xl, yl = points[l]
                    # 以点 i 为矩形的一个顶点，邻边指向 k 和 l
                    # 矩形面积 = |(k - i) × (l - i)|
                    area = abs((xk - x1) * (yl - y1) - (yk - y1) * (xl - x1))
                    if area < res:
                        res = area

                diag_map[key].append((i, j))

        return res if res < INF else 0.0
