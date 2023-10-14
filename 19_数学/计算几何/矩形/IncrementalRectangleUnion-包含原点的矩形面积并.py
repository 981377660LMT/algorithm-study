"""
!`包含原点的`矩形合并/矩形面积并
维护一些矩形的面积并,支持查询面积和添加矩形。
注意矩形的范围在第一象限，且矩形左下角坐标为(0,0),右上角坐标为(x,y)。

三维情形:
扫描z轴,维护一个二维的矩形面积并。
https://maspypy.github.io/library/other/cuboid_union_volume.hpp
https://maspypy.github.io/library/ds/incremental_rectangle_union.hpp
https://kmjp.hatenablog.jp/entry/2017/06/20/0900
https://ei1333.github.io/library/structure/others/union-rectangle.hpp
"""
# IncrementalRectangleUnion-包含原点的矩形面积并/RectangleUnion/IncrementalRectangleUnion


from typing import List, Tuple
from sortedcontainers import SortedList

INF = int(1e18)


class IncrementalRectangleUnionRange:
    __slots__ = ("_ur", "_ul", "_dr", "_dl")

    def __init__(self) -> None:
        self._ur = IncrementalRectangleUnion()
        self._ul = IncrementalRectangleUnion()
        self._dr = IncrementalRectangleUnion()
        self._dl = IncrementalRectangleUnion()

    def add(self, x1: int, x2: int, y1: int, y2: int) -> None:
        """Add [x1, x2] * [y1, y2].
        x1<=0<=x2.
        y1<=0<=y2.
        """
        assert x1 <= 0 <= x2, y1 <= 0 <= y2
        self._ur.add(x2, y2)
        self._ul.add(-x1, y2)
        self._dr.add(x2, -y1)
        self._dl.add(-x1, -y1)

    def query(self) -> int:
        """Return the sum of all rectangles."""
        return self._ur.sum + self._ul.sum + self._dr.sum + self._dl.sum


class IncrementalRectangleUnion:
    __slots__ = ("_sl", "sum")

    def __init__(self) -> None:
        self._sl = SortedList([(0, INF), (INF, 0)])
        self.sum = 0

    def add(self, x: int, y: int) -> None:
        """Add [0, x] * [0, y]."""
        pos = self._sl.bisect_left((x, -INF))
        item = self._sl[pos]
        if item[1] >= y:
            return
        nextY = item[1]
        pos -= 1
        pre = self._sl[pos]
        while pre[1] <= y:
            x1, y1 = pre
            del self._sl[pos]
            pos -= 1
            pre = self._sl[pos]
            self.sum -= (x1 - pre[0]) * (y1 - nextY)
        self.sum += (x - self._sl[pos][0]) * (y - nextY)
        pos = self._sl.bisect_left((x, -INF))
        if self._sl[pos][0] == x:
            self._sl.pop(pos)
        self._sl.add((x, y))

    def query(self) -> int:
        """Return the sum of all rectangles."""
        return self.sum


def cuboidUnionVolumn(points: List[Tuple[int, int, int]]) -> int:
    """求出所有左下角为原点的立方体的体积并.

    Args:
        points (List[Tuple[int,int,int]]): [x, y, z] 每个点的坐标，表示一个[0, x] * [0, y] * [0, z]的立方体.

    Returns:
        int: 所有立方体的体积并.
    """
    points = sorted(points, key=lambda x: x[2], reverse=True)
    preZ = INF
    res = 0
    area = 0
    manager = IncrementalRectangleUnion()
    for x, y, z in points:
        res += (preZ - z) * area
        manager.add(x, y)
        area = manager.query()
        preZ = z
    res += preZ * area
    return res


if __name__ == "__main__":
    ur = IncrementalRectangleUnionRange()

    ur.add(-2, 1, -2, 1)
    print(ur.query())
    ur.add(-1, 2, -1, 2)
    print(ur.query())

    # https://yukicoder.me/problems/2577
    n = int(input())
    adds = [tuple(map(int, input().split())) for _ in range(n)]  # (x1,y1,x2,y2)
    UR = IncrementalRectangleUnionRange()
    pre = 0
    for x1, y1, x2, y2 in adds:
        UR.add(x1, x2, y1, y2)
        cur = UR.query()
        print(cur - pre)
        pre = cur

    class Solution:
        def rectangleArea(self, rectangles: List[List[int]]) -> int:
            ur = IncrementalRectangleUnionRange()
            for x1, y1, x2, y2 in rectangles:
                ur.add(x1, x2, y1, y2)
            return ur.query() % int(1e9 + 7)
