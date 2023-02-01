# 二维平面，给定一个凸多边形，以及在凸多边形外的两个点s,t。
# 问从点 s到点 t的最短距离，期间不能经过凸多边形`内部`。

# n<=1e5
# 凸多边形的任意三点不共线,给出的点为逆时针顺序


# 解:
# 将n+2个点求一遍凸包
# 如果s或t不在凸包上,则答案为s到t的距离
# 如果s和t都在凸包上,则答案为凸包上顺时针/逆时针的两条路径中较短的一条

from typing import List, Tuple

Point = Tuple[int, int]


def dontSwim(convexPoly: List[Point], start: Point, target: Point) -> float:
    convexPoly += [start, target]
    hull = calConvexHull(convexPoly)
    sid, tid = -1, -1
    for i, p in enumerate(hull):
        if p == start:
            sid = i
        if p == target:
            tid = i
    if sid == -1 or tid == -1:
        return calDist(start, target)

    if sid > tid:
        sid, tid = tid, sid
    path1 = hull[sid : tid + 1]
    path2 = hull[tid:] + hull[: sid + 1]
    dist1 = sum(calDist(path1[i], path1[i + 1]) for i in range(len(path1) - 1))
    dist2 = sum(calDist(path2[i], path2[i + 1]) for i in range(len(path2) - 1))
    return min(dist1, dist2)


def calDist(A: Point, B: Point) -> float:
    """计算两点间距离"""
    dx, dy = A[0] - B[0], A[1] - B[1]
    return (dx * dx + dy * dy) ** 0.5


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
    n = int(input())
    convextPoly = list(tuple(map(int, input().split())) for _ in range(n))
    start = tuple(map(int, input().split()))
    target = tuple(map(int, input().split()))
    print(dontSwim(convextPoly, start, target))
