# from scipy.spatial import ConvexHull
from typing import List

Point = List[int]


def cross(a: Point, b: Point, c: Point) -> int:
    """计算AB与AC的叉乘"""
    return (b[0] - a[0]) * (c[1] - a[1]) - (b[1] - a[1]) * (c[0] - a[0])


def convexHull(points: List[Point]) -> List[Point]:
    """Andrew 算法 nlogn 求凸包"""
    points = sorted(points)
    stack = []

    # 寻找凸壳的下半部分
    for i in range(len(points)):
        while len(stack) >= 2 and cross(stack[-2], stack[-1], points[i]) < 0:
            stack.pop()
        stack.append(tuple(points[i]))

    # 寻找凸壳的上半部分
    for i in range(len(points) - 1, -1, -1):
        while len(stack) >= 2 and cross(stack[-2], stack[-1], points[i]) < 0:
            stack.pop()
        stack.append(tuple(points[i]))

    return list(set(stack))


if __name__ == "__main__":
    expected = [(2, 4), (1, 1), (2, 0), (4, 2), (3, 3)]
    print(convexHull([[1, 1], [2, 2], [2, 0], [2, 4], [3, 3], [4, 2]]))
