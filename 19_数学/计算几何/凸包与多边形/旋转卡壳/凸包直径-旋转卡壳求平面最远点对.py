# https://github.com/EndlessCheng/codeforces-go/blob/a707a2c9af5063a42fae95bcd38a0be21ea600cc/copypasta/geometry.go#L859
# 凸包直径-旋转卡壳求平面最远点对
# https://www.luogu.com.cn/problem/P1452
# n<=5e4


from typing import List, Tuple


def det(a: Tuple[int, int], b: Tuple[int, int]) -> int:
    return a[0] * b[1] - a[1] * b[0]


def convexHull(points: List[Tuple[int, int]]) -> List[Tuple[int, int]]:
    """葛立恒扫描法 Graham's scan 求凸包"""
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


def rotatingCalipers(points: List[Tuple[int, int]]) -> Tuple[Tuple[int, int], Tuple[int, int]]:
    ch = convexHull(points)
    n = len(ch)
    if n == 2:
        return ch[0], ch[1]
    i, j = 0, 0
    for k, p in enumerate(ch):
        if not ch[i] < p:
            i = k
        if ch[j] < p:
            j = k
    maxD2 = 0
    i0, j0 = i, j
    p1, p2 = (), ()
    while i != j0 or j != i0:
        dx, dy = ch[i][0] - ch[j][0], ch[i][1] - ch[j][1]
        d2 = dx * dx + dy * dy
        if d2 > maxD2:
            maxD2 = d2
            p1, p2 = ch[i], ch[j]
        if (
            det(
                (ch[(i + 1) % n][0] - ch[i][0], ch[(i + 1) % n][1] - ch[i][1]),
                (ch[(j + 1) % n][0] - ch[j][0], ch[(j + 1) % n][1] - ch[j][1]),
            )
            < 0
        ):
            i = (i + 1) % n
        else:
            j = (j + 1) % n
    return p1, p2  # type: ignore


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    points = list(tuple(map(int, input().split())) for _ in range(n))
    p1, p2 = rotatingCalipers(points)
    print((p1[0] - p2[0]) ** 2 + (p1[1] - p2[1]) ** 2)
