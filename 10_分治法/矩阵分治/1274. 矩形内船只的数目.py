class Point(object):
    def __init__(self, x: int, y: int):
        self.x = x
        self.y = y


# 当且仅当这两个点所表示的矩形区域（包含边界）内至少有一艘船时，这个函数才返回 true ，否则返回 false 。
class Sea(object):
    def hasShips(self, topRight: 'Point', bottomLeft: 'Point') -> bool:
        ...


# 给你矩形的右上角 topRight 和左下角 bottomLeft 的坐标，请你返回此矩形内船只的数目。题目保证矩形内 至多只有 10 艘船。
class Solution(object):
    def countShips(self, sea: 'Sea', topRight: 'Point', bottomLeft: 'Point') -> int:
        x1, y1 = topRight.x, topRight.y
        x2, y2 = bottomLeft.x, bottomLeft.y

        if x1 < x2 or y1 < y2 or not sea.hasShips(topRight, bottomLeft):
            return 0
        if (x1, y1) == (x2, y2):
            return 1

        xmid, ymid = (x1 + x2) // 2, (y1 + y2) // 2
        return (
            self.countShips(sea, Point(xmid, ymid), Point(x2, y2))
            + self.countShips(sea, Point(xmid, y1), Point(x2, ymid + 1))
            + self.countShips(sea, Point(x1, ymid), Point(xmid + 1, y2))
            + self.countShips(sea, Point(x1, y1), Point(xmid + 1, ymid + 1))
        )

