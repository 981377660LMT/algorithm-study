# 三分法求单峰函数的极值点
# https://oi-wiki.org/basic/binary/#_10
from math import dist
from typing import List

EPS = 1e-6


class Solution:
    def getMinDistSum(self, positions: List[List[int]]) -> float:
        """三分求一维凸函数的最值，二维就要对横纵坐标三分两次"""

        def calDistSum(centerX: float, centerY: float) -> float:
            return sum(dist((x, y), (centerX, centerY)) for x, y in positions)

        def trisectX(centerY: float) -> float:
            leftX, rightX = 0.0, 100.0
            while leftX <= rightX:
                diff = rightX - leftX
                mid1 = leftX + diff / 3
                mid2 = leftX + 2 * diff / 3
                if calDistSum(mid1, centerY) < calDistSum(mid2, centerY):
                    rightX = mid2 - EPS
                else:
                    leftX = mid1 + EPS
            return calDistSum(leftX, centerY)

        leftY, rightY = 0.0, 100.0
        while leftY <= rightY:
            diff = rightY - leftY
            mid1 = leftY + diff / 3
            mid2 = leftY + 2 * diff / 3
            if trisectX(mid1) < trisectX(mid2):
                rightY = mid2 - EPS
            else:
                leftY = mid1 + EPS
        return trisectX(leftY)
