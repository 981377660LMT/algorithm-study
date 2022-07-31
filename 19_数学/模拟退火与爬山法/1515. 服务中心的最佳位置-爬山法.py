# 所求的最佳位置即为几何中位数（Geometric Median (opens new window)）(费马点)
# 模拟退火是随机选择临近状态，根据温度概率性选择是否转移
from math import dist, sqrt
from typing import List
from scipy.optimize import minimize

EPS = 1e-6
# 请你为服务中心选址，该位置的坐标 [xcentre, ycentre] 需要使下面的公式取到最小值：
# 与真实值误差在 10-5之内的答案将被视作正确答案。

# 0 <= xi, yi <= 100
class Solution:
    def getMinDistSum1(self, positions: List[List[int]]) -> float:
        # """https://leetcode-cn.com/problems/best-position-for-a-service-centre/solution/5463-fu-wu-zhong-xin-de-zui-jia-wei-zhi-by-tuotuol/"""
        # !每轮按比例缩短步伐 => 爬山法
        dist = lambda p: sum(sqrt((p[0] - i) ** 2 + (p[1] - j) ** 2) for i, j in positions)
        dirs = [(0, 50), (50, 0), (0, -50), (-50, 0)]
        x, y = 50, 50
        for _ in range(90):
            x, y = min(((x + i, y + j) for i, j in dirs), key=dist)
            dirs = [(i / 1.2, j / 1.2) for i, j in dirs]
        return dist((x, y))

    def getMinDistSum2(self, positions: List[List[int]]) -> float:
        # https://www.runoob.com/scipy/scipy-optimize.html
        # https://leetcode.cn/problems/best-position-for-a-service-centre/solution/by-freeyourmind-k558/
        def getSum(target: List[int]):
            return sum(
                [sqrt((target[0] - x1) ** 2 + (target[1] - y1) ** 2) for x1, y1 in positions]
            )

        return minimize(getSum, [50, 50]).fun  # fun is fun

    def getMinDistSum(self, positions: List[List[int]]) -> float:
        """三分求一维凸函数的极值，二维就要对横纵坐标三分两次"""

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
