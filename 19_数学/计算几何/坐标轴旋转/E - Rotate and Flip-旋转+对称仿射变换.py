# E - Rotate and Flip-旋转+对称 (仿射变换 affine transformation)
# 给定n个点，每个点有一个编号，第i个点的坐标为(xi,yi)。有m个操作，操作有4种，输入格式和操作内容如下所示。
# 1：将所有点沿着原点为中心的顺时针旋转90度
# 2：将所有点沿着原点为中心的逆时针旋转90度
# 3 p：将所有点沿着直线x=p对称
# 4 p：将所有点沿着直线y=p对称

# 给定q个询问，每个询问有两个整数ai,bi，表示在执行ai次操作后，第bi个点的坐标是多少。
# !n,m,q<=2e5  所有坐标绝对值<=1e9
# 0<=ai<=m 1<=bi<=n

# !每次只需要求出 (0,0),(1,0),(0,1) 移动到了哪里，然后根据这三个点的坐标，可以求出其他点的坐标
# !每个仿射变换(a, b, c, d, e, f)都满足 (x, y) => (ax + by + c, dx + ey + f)
# !每一个仿射变换都可以写成一个矩阵A，矩阵乘法可以求出变换后的坐标 Ap1 => p2
# https://atcoder.jp/contests/abc189/editorial/539

from typing import List, Tuple


Point = Tuple[int, int]


def rotate90degClockWisely(point: "Point") -> "Point":
    """(x, y)顺时针旋转90度后的坐标"""
    return point[1], -point[0]


def rotate90degCounterClockWisely(point: "Point") -> "Point":
    """(x, y)逆时针旋转90度后的坐标"""
    return -point[1], point[0]


def flipX(point: "Point", p: int) -> "Point":
    """(x, y)关于直线x=p对称后的坐标"""
    return 2 * p - point[0], point[1]


def flipY(point: "Point", p: int) -> "Point":
    """(x, y)关于直线y=p对称后的坐标"""
    return point[0], 2 * p - point[1]


def affineTransformation(point: "Point", reference: Tuple["Point", "Point", "Point"]) -> "Point":
    """
    求(x, y)经过坐标系仿射变换后的坐标

    - Args:
        point: 原坐标
        reference: (0,0),(1,0),(0,1)经过仿射变换后的坐标

    - Returns:
        经过仿射变换后的坐标

    - Notes:
        每个仿射变换(a, b, c, d, e, f)都满足 (x, y) => (ax + by + c, dx + ey + f)
    """
    c, f = reference[0]
    a = reference[1][0] - c
    b = reference[2][0] - c
    d = reference[1][1] - f
    e = reference[2][1] - f
    return a * point[0] + b * point[1] + c, d * point[0] + e * point[1] + f


def rotateAndFlip(
    points: List["Point"], operations: List[List[int]], queries: List[Tuple[int, int]]
) -> List["Point"]:
    pos: List[Tuple["Point", "Point", "Point"]] = [
        ((0, 0), (1, 0), (0, 1))
    ]  # 每轮操作后，(0,0),(1,0),(0,1)的坐标
    for op, *args in operations:
        if op == 1:
            pos.append(tuple(rotate90degClockWisely(point) for point in pos[-1]))
        elif op == 2:
            pos.append(tuple(rotate90degCounterClockWisely(point) for point in pos[-1]))
        elif op == 3:
            p = args[0]
            pos.append(tuple(flipX(point, p) for point in pos[-1]))
        elif op == 4:
            p = args[0]
            pos.append(tuple(flipY(point, p) for point in pos[-1]))

    res = []
    for i, j in queries:
        reference = pos[i]
        point = points[j - 1]
        res.append(affineTransformation(point, reference))
    return res


if __name__ == "__main__":
    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    m = int(input())
    operations = [list(map(int, input().split())) for _ in range(m)]
    q = int(input())
    queries = [tuple(map(int, input().split())) for _ in range(q)]
    res = rotateAndFlip(points, operations, queries)
    for x, y in res:
        print(x, y)
