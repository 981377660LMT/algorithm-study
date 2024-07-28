# 100347. 判断矩形的两个角落是否可达
# https://leetcode.cn/problems/check-if-the-rectangle-corner-is-reachable/description/
# 给你两个正整数 X 和 Y 和一个二维整数数组 circles ，其中 circles[i] = [xi, yi, ri] 表示一个圆心在 (xi, yi) 半径为 ri 的圆。
# 坐标平面内有一个左下角在原点，右上角在 (X, Y) 的矩形。你需要判断是否存在一条从左下角到右上角的路径满足：路径 完全 在矩形内部，不会 触碰或者经过 任何 圆的内部和边界，同时 只 在起点和终点接触到矩形。
# 如果存在这样的路径，请你返回 true ，否则返回 false 。
# !圆心一定在矩形内
#
# 是否可达 => 并查集

from UnionFindArraySimple import UnionFindArraySimple

from typing import List


def isCircleIntersectCircle(cx1: int, cy1: int, r1: int, cx2: int, cy2: int, r2: int) -> bool:
    """两个圆是否有交点."""
    dx, dy = cx1 - cx2, cy1 - cy2
    return dx * dx + dy * dy <= (r1 + r2) * (r1 + r2)


def isCircleIntersectX(cx: int, cy: int, r: int, xa: int) -> bool:
    """圆和直线x=xa是否有交点."""
    return cx - r <= xa <= cx + r


def isCircleIntersectY(cx: int, cy: int, r: int, ya: int) -> bool:
    """圆和直线y=ya是否有交点."""
    return cy - r <= ya <= cy + r


class Solution:
    def canReachCorner(self, X: int, Y: int, circles: List[List[int]]) -> bool:
        n = len(circles)
        LEFT_UP, RIGHT_DOWN = n, n + 1

        # !并查集维护相交圆的分组
        uf = UnionFindArraySimple(n + 2)
        for i, (x1, y1, r1) in enumerate(circles):
            if isCircleIntersectX(x1, y1, r1, 0) or isCircleIntersectY(x1, y1, r1, Y):
                uf.union(i, LEFT_UP)
            if isCircleIntersectX(x1, y1, r1, X) or isCircleIntersectY(x1, y1, r1, 0):
                uf.union(i, RIGHT_DOWN)

            for j, (x2, y2, r2) in enumerate(circles[:i]):
                if isCircleIntersectCircle(x1, y1, r1, x2, y2, r2):
                    uf.union(i, j)

            if uf.find(LEFT_UP) == uf.find(RIGHT_DOWN):
                return False

        return True
