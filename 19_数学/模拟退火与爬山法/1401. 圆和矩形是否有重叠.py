# 1401. 圆和矩形是否有重叠
# https://leetcode.cn/problems/circle-and-rectangle-overlapping/solution/xi-jie-tai-duo-san-fen-tao-san-fen-bing-o6f3l/
#
# 三分套三分
# f(x,y)=(x−xCenter)^2+(y−yCenter)^2−radius^2
# x固定时，f(y)是一个单峰函数，可以用三分法求最值
# 只要f(x,y)最小值<=0,就说明圆和矩形有重叠部分
#
# 给你一个以 (radius, xCenter, yCenter) 表示的圆和一个与坐标轴平行的矩形 (x1, y1, x2, y2) ，
# 其中 (x1, y1) 是矩形左下角的坐标，而 (x2, y2) 是右上角的坐标。
# 如果圆和矩形有重叠的部分，请你返回 true ，否则返回 false 。
# 换句话说，请你检测是否 存在 点 (xi, yi) ，它既在圆上也在矩形上（两者都包括点落在边界上的情况）。


class Solution:
    def checkOverlap(
        self, radius: int, xCenter: int, yCenter: int, x1: int, y1: int, x2: int, y2: int
    ) -> bool:
        def calDistSum(x: int, y: int) -> int:
            return (x - xCenter) * (x - xCenter) + (y - yCenter) * (y - yCenter) - radius * radius

        def trisectY(curX: int) -> int:
            """固定x,三分y求f(x,y)最小值"""
            leftY, rightY = y1, y2
            while rightY - leftY >= 3:
                diff = rightY - leftY
                mid1 = leftY + diff // 3
                mid2 = leftY + 2 * diff // 3
                if calDistSum(curX, mid1) > calDistSum(curX, mid2):
                    leftY = mid1
                else:
                    rightY = mid2
            min_ = calDistSum(curX, leftY)
            leftY += 1
            while leftY <= rightY:
                cand = calDistSum(curX, leftY)
                min_ = cand if cand < min_ else min_
                leftY += 1
            return min_

        leftX, rightX = x1, x2
        while rightX - leftX >= 3:
            diff = rightX - leftX
            mid1 = leftX + diff // 3
            mid2 = leftX + 2 * diff // 3
            if trisectY(mid1) > trisectY(mid2):
                leftX = mid1
            else:
                rightX = mid2
        while leftX <= rightX:
            if trisectY(leftX) <= 0:
                return True
            leftX += 1

        return False
