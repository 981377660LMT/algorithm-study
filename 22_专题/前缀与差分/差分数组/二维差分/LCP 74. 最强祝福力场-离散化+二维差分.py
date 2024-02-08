# LCP 74. 最强祝福力场-离散化+二维差分
# https://leetcode.cn/problems/xepqZ5/
# forceField[i] = [x,y,side] 表示第 i 片力场将覆盖以坐标 (x,y) 为中心，边长为 side 的正方形区域。
# !若任意一点的 力场强度 等于覆盖该点的力场数量，请求出在这片地带中 力场强度 最强处的 力场强度。

# !统计所有左下和右上坐标，由于会出现 0.5可以将坐标乘 2。
# O(n^2)


from typing import List
from Diff2D import Diff2D


class Solution:
    def fieldOfGreatestBlessing(self, forceField: List[List[int]]) -> int:
        # 离散化
        allX, allY = set(), set()
        for x, y, side in forceField:
            allX.add(2 * x - side)
            allX.add(2 * x + side)
            allY.add(2 * y - side)
            allY.add(2 * y + side)
        sortedX = sorted(allX)
        sortedY = sorted(allY)
        rankX = {x: i for i, x in enumerate(sortedX)}
        rankY = {y: i for i, y in enumerate(sortedY)}

        # 二维差分
        row, col = len(sortedX), len(sortedY)
        diffMatrix = Diff2D([[0] * col for _ in range(row)])
        for x, y, side in forceField:
            r1, c1 = rankX[2 * x - side], rankY[2 * y - side]
            r2, c2 = rankX[2 * x + side], rankY[2 * y + side]
            diffMatrix.add(r1, c1, r2, c2, 1)

        res = 0
        for i in range(row):
            for j in range(col):
                res = max(res, diffMatrix.query(i, j))
        return res
