"""
矩形合并/矩形的面积并
维护一些矩形的面积并,支持查询面积和添加矩形。
注意矩形的范围在第一象限，且矩形左下角坐标为(0,0),右上角坐标为(x,y)。
"""

from sortedcontainers import SortedList


INF = int(1e18)


class UnionRectangle:
    __slots__ = ("_data", "_sum")

    def __init__(self) -> None:
        self._data = SortedList([(0, INF), (INF, 0)])
        self._sum = 0

    def addPoint(self, x: int, y: int) -> None:
        """Add [0, x] * [0, y]."""
        pos = self._data.bisect_left((x, -INF))
        if self._data[pos][1] >= y:
            return
        nextY = self._data[pos][1]
        pos -= 1
        while self._data[pos][1] <= y:
            x1, y1 = self._data[pos]
            del self._data[pos]
            pos -= 1
            self._sum -= (x1 - self._data[pos][0]) * (y1 - nextY)
        self._sum += (x - self._data[pos][0]) * (y - nextY)
        self._data.add((x, y))

    def query(self) -> int:
        """Return the sum of all rectangles."""
        return self._sum


if __name__ == "__main__":
    points = [(2, 2), (3, 4), (1, 7)]
    ur = UnionRectangle()
    for x, y in points:
        ur.addPoint(x, y)
        print(ur.query())
