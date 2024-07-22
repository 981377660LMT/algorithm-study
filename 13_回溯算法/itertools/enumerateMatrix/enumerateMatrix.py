# 二维矩阵遍历工具函数(NOT VERIFIED)
# 1. 顺时针遍历螺旋矩阵
# 2. 顺时针遍历矩阵从外向内的第 d 圈（保证不自交）
# 3. 之字遍历矩阵
# 4. 遍历以 (ox, oy) 为中心的曼哈顿距离为 dist 的边界上的格点
# 5. 从上到下，遍历以 (ox, oy) 为中心的曼哈顿距离 <= dist 的所有格点
# 6. 曼哈顿圈序遍历
# 7. 遍历以 (ox, oy) 为中心的切比雪夫距离为 dist 的边界上的格点
# 8. 遍历矩阵的边界


from typing import Generator, List, Tuple


def enumerateMatrixSpirally(row: int, col: int) -> Generator[Tuple[int, int], None, None]:
    """顺时针遍历row行col列的螺旋矩阵(SpiralMatrix)."""
    dir4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]  # 右下左上
    mat = [[-1 for _ in range(col)] for _ in range(row)]
    r, c, di = 0, 0, 0
    for id in range(row * col):
        yield r, c
        mat[r][c] = id
        d = dir4[di]
        nx, ny = r + d[0], c + d[1]
        if nx < 0 or nx >= row or ny < 0 or ny >= col or mat[nx][ny] != -1:
            di = (di + 1) & 3
            d = dir4[di]
        r += d[0]
        c += d[1]


def enumerateMatrixAroundByDepth(
    mat: List[List[int]], depth: int
) -> Generator[Tuple[int, int], None, None]:
    """顺时针遍历矩阵从外向内的第 d 圈（保证不自交）.
    一共 (row+col-d*4-2)*2 个元素.
    """
    row, col = len(mat), len(mat[0])
    for j in range(depth, col - depth):  # →
        yield depth, j
    for i in range(depth + 1, row - depth):  # ↓
        yield i, col - 1 - depth
    for j in range(col - depth - 2, depth - 1, -1):  # ←
        yield row - 1 - depth, j
    for i in range(row - depth - 2, depth, -1):  # ↑
        yield i, depth


def enumerateMatrixZigZag(row: int, col: int) -> Generator[Tuple[int, int], None, None]:
    """获取之字遍历的所有坐标.
    例如 3 行 3 列的矩阵：
    0 -> 1 -> 2
              ↓
    3 <- 4 <- 5
    ↓
    6 -> 7 -> 8
    """
    for i in range(row):
        for j in range(col):
            yield i, j
        i += 1
        if i == row:
            break
        for j in range(col - 1, -1, -1):
            yield i, j


def enumerateMatrixAroundByManhattan(
    n: int, m: int, ox: int, oy: int, dist: int
) -> Generator[Tuple[int, int], None, None]:
    """遍历以 (ox, oy) 为中心的曼哈顿距离为 dist 的边界上的格点.
    从最右顶点出发，逆时针移动.
    """
    if dist == 0:
        yield ox, oy
        return
    dir4r = [(-1, 1), (-1, -1), (1, -1), (1, 1)]  # 逆时针
    x, y = ox + dist, oy  # 从最右顶点出发，逆时针移动
    for d in dir4r:
        for _ in range(dist):
            if 0 <= x < n and 0 <= y < m:
                yield x, y
            x += d[0]
            y += d[1]


def enumerateMatrixInnerByManhattan(
    n: int, m: int, ox: int, oy: int, dist: int
) -> Generator[Tuple[int, int], None, None]:
    """从上到下，遍历以 (ox, oy) 为中心的曼哈顿距离 <= dist 的所有格点."""
    for i in range(max2(ox - dist, 0), min2(ox + dist + 1, n)):
        d = dist - abs(ox - i)
        for j in range(max2(oy - d, 0), min2(oy + d + 1, m)):
            if i == ox and j == oy:
                continue
            yield i, j


def enumerateMatrixByManhattan(
    n: int, m: int, ox: int, oy: int
) -> Generator[Tuple[int, int], None, None]:
    """曼哈顿圈序遍历.
    从最右顶点出发，逆时针移动.
    """
    yield ox, oy
    dir4r = [(-1, 1), (-1, -1), (1, -1), (1, 1)]  # 逆时针
    maxDist = max2(ox, n - 1 - ox) + max2(oy, m - 1 - oy)
    for dis in range(1, maxDist + 1):
        x, y = ox + dis, oy  # 从最右顶点出发，逆时针移动
        for d in dir4r:
            for _ in range(dis):
                if 0 <= x < n and 0 <= y < m:
                    yield x, y
                x += d[0]
                y += d[1]


def enumerateMatrixAroundByChebyshev(
    n: int, m: int, ox: int, oy: int, dist: int
) -> Generator[Tuple[int, int], None, None]:
    """遍历以 (ox, oy) 为中心的切比雪夫距离为 dist 的边界上的格点.
    从最右顶点出发，逆时针移动.
    """
    # 上下
    for x in [ox - dist, ox + dist]:
        if 0 <= x < n:
            for y in range(max2(oy - dist, 0), min2(oy + dist + 1, m)):
                yield x, y
    # 左右（注意四角已经被上面的循环枚举到了）
    for y in [oy - dist, oy + dist]:
        if 0 <= y < m:
            for x in range(max2(ox - dist, 0) + 1, min2(ox + dist, n - 1)):
                yield x, y


def enumerateMatrixBorder(
    x0: int, x1: int, y0: int, y1: int
) -> Generator[Tuple[int, int], None, None]:
    """遍历矩阵的边界.
    x0<=x<=x1, y0<=y<=y1.
    """
    if y0 == y1:
        for i in range(x0, x1 + 1):
            yield i, y0
        return
    for i in range(x0, x1 + 1):
        j = y0
        while j <= y1:
            yield i, j
            if i == x0 or i == x1:
                j += 1
            else:
                j += y1 - y0


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b
