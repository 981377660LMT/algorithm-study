from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个正整数 X 和 Y 和一个二维整数数组 circles ，其中 circles[i] = [xi, yi, ri] 表示一个圆心在 (xi, yi) 半径为 ri 的圆。

# 坐标平面内有一个左下角在原点，右上角在 (X, Y) 的矩形。你需要判断是否存在一条从左下角到右上角的路径满足：路径 完全 在矩形内部，不会 触碰或者经过 任何 圆的内部和边界，同时 只 在起点和终点接触到矩形。


# 如果存在这样的路径，请你返回 true ，否则返回 false 。
# !圆心一定在矩形内
#
# 是否可达 => 并查集


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(self, key1: int, key2: int) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


def isPointInCircle(x: int, y: int, cx: int, cy: int, r: int) -> bool:
    dx, dy = x - cx, y - cy
    return dx * dx + dy * dy <= r * r


def isCircleIntersectCircle(cx1: int, cy1: int, r1: int, cx2: int, cy2: int, r2: int) -> bool:
    dx, dy = cx1 - cx2, cy1 - cy2
    return dx * dx + dy * dy <= (r1 + r2) * (r1 + r2)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def canReachCorner(self, X: int, Y: int, circles: List[List[int]]) -> bool:
        # !起点或者终点不能在圆上、圆内
        for x, y, r in circles:
            if isPointInCircle(0, 0, x, y, r) or isPointInCircle(X, Y, x, y, r):
                return False

        # !并查集维护相交圆的分组
        uf = UnionFindArraySimple(len(circles))
        for i in range(len(circles)):
            x1, y1, r1 = circles[i]
            for j in range(i + 1, len(circles)):
                x2, y2, r2 = circles[j]
                if isCircleIntersectCircle(x1, y1, r1, x2, y2, r2):
                    uf.union(i, j)

        groups = defaultdict(list)
        for i in range(len(circles)):
            groups[uf.find(i)].append(i)

        for group in groups.values():
            minX, minY, maxX, maxY = INF, INF, -INF, -INF
            for i in group:
                x, y, r = circles[i]
                minX, minY = min2(minX, x - r), min2(minY, y - r)
                maxX, maxY = max2(maxX, x + r), max2(maxY, y + r)
            xCovered = minX <= 0 and maxX >= X
            if xCovered:
                return False
            yCovered = minY <= 0 and maxY >= Y
            if yCovered:
                return False
            startClosed = minX <= 0 and minY <= 0
            if startClosed:
                return False
            targetClosed = maxX >= X and maxY >= Y
            if targetClosed:
                return False
        return True


# X = 3, Y = 3, circles = [[2,1,1],[1,2,1]]

print(Solution().canReachCorner(3, 3, [[2, 1, 1], [1, 2, 1]]))
# 5
# 8
# [[4,7,1]]
print(Solution().canReachCorner(5, 8, [[4, 7, 1]]))
