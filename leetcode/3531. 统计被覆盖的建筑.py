# 3531. 统计被覆盖的建筑
# https://leetcode.cn/problems/count-covered-buildings/description/
# 给你一个正整数 n，表示一个 n x n 的城市，同时给定一个二维数组 buildings，其中 buildings[i] = [x, y] 表示位于坐标 [x, y] 的一个 唯一 建筑。
# 如果一个建筑在四个方向（左、右、上、下）中每个方向上都至少存在一个建筑，则称该建筑 被覆盖 。
# 返回 被覆盖 的建筑数量。


from typing import List


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def countCoveredBuildings(self, n: int, buildings: List[List[int]]) -> int:
        xToMinY = [n + 1] * (n + 1)
        xToMaxY = [-1] * (n + 1)
        yToMinX = [n + 1] * (n + 1)
        yToMaxX = [-1] * (n + 1)

        for x, y in buildings:
            xToMinY[x] = min2(xToMinY[x], y)
            xToMaxY[x] = max2(xToMaxY[x], y)
            yToMinX[y] = min2(yToMinX[y], x)
            yToMaxX[y] = max2(yToMaxX[y], x)

        res = 0
        for x, y in buildings:
            if xToMinY[x] < y < xToMaxY[x] and yToMinX[y] < x < yToMaxX[y]:
                res += 1
        return res
