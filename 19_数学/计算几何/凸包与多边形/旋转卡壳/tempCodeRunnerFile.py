
from typing import List, Tuple


def maxArea(points: List[Tuple[float, float]]) -> float:
    """最大四边形面积"""
    ch = convexHull(points)
    res = 0
    n = len(points)
    for i in range(n):
        a, b = i, (i + 1) % n
        for j in range(i + 1, n):
            while cross(ch[i], ch[j], ch[(a + 1) % n]) < cross(ch[i], ch[j], ch[a]):
                a = (a + 1) % n
            while cross(ch[i], ch[j], ch[(b + 1) % n]) > cross(ch[i], ch[j], ch[b]):
                b = (b + 1) % n
            cand = -cross(ch[i], ch[j], ch[a]) + cross(ch[i], ch[j], ch[b])
            res = cand if cand > res else res
    return res / 2


def cross(a: Tuple[float, float], b: Tuple[float, float], c: Tuple[float, float]) -> float:
    return (b[0] - a[0]) * (c[1] - a[1]) - (c[0] - a[0]) * (b[1] - a[1])


def det(a: Tuple[float, float], b: Tuple[float, float]) -> float:
    return a[0] * b[1] - a[1] * b[0]


def convexHull(points: List[Tuple[float, float]]) -> List[Tuple[float, float]]:
    points = sorted(points)
    res = []
    for p in points:
        while (
            len(res) > 1
            and det(
                (res[-1][0] - res[-2][0], res[-1][1] - res[-2][1]),
                (p[0] - res[-1][0], p[1] - res[-1][1]),
            )
            <= 0
        ):
            res.pop()
        res.append(p)
    sz = len(res)
    for i in range(len(points) - 2, -1, -1):
        p = points[i]
        while (
            len(res) > sz
            and det(
                (res[-1][0] - res[-2][0], res[-1][1] - res[-2][1]),
                (p[0] - res[-1][0], p[1] - res[-1][1]),
            )
            <= 0
        ):
            res.pop()
        res.append(p)
    res.pop()  # 如果需要首尾相同则去掉这行
    return res


if __name__ == "__main__":
    n = int(input())
    points = list(tuple(map(float, input().split())) for _ in range(n))
    res = maxArea(points)
    print(f"{res:.3f}")
