# from scipy.spatial import ConvexHull
from typing import List

Point = List[int]


def calCross(A: Point, B: Point, C: Point) -> int:
    """计算AB与AC的叉乘"""

    AB = [B[0] - A[0], B[1] - A[1]]
    AC = [C[0] - A[0], C[1] - A[1]]
    return AB[0] * AC[1] - AB[1] * AC[0]


def calConvexHull(points: List[Point]) -> List[Point]:
    """Andrew 算法 nlogn 求凸包"""
    if len(points) <= 3:
        return points

    points = sorted(points)
    stack = []

    # 寻找凸壳的下半部分
    for i in range(len(points)):
        while len(stack) >= 2 and calCross(stack[-2], stack[-1], points[i]) < 0:
            stack.pop()
        stack.append(tuple(points[i]))

    # 寻找凸壳的上半部分
    for i in range(len(points) - 1, -1, -1):
        while len(stack) >= 2 and calCross(stack[-2], stack[-1], points[i]) < 0:
            stack.pop()
        stack.append(tuple(points[i]))
    return list(set(stack))


if __name__ == "__main__":
    expected = [(2, 4), (1, 1), (2, 0), (4, 2), (3, 3)]
    assert calConvexHull([[1, 1], [2, 2], [2, 0], [2, 4], [3, 3], [4, 2]]) == expected
