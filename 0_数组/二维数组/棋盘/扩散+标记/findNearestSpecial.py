# getNearestSpecial/findNearestSpecial

from typing import Callable, List, Tuple


BoundingRect = Tuple[int, int, int, int]  # (top,bottom,left,right)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


def findNearestSpecial(
    row: int, col: int, isSpecialFn: Callable[[Tuple[int, int]], bool]
) -> List[List[BoundingRect]]:
    """
    给定一个row*col的矩阵, 对每个点, 找到 top,bottom,left,right 四个方向上最近的特殊点(包含自身).
    如果top不存在, 则top=-1;
    如果bottom不存在, 则bottom=row;
    如果left不存在, 则left=-1;
    如果right不存在, 则right=col.
    """
    isSpecial = [[isSpecialFn((r, c)) for c in range(col)] for r in range(row)]
    leftRes = [[-1] * col for _ in range(row)]
    rightRes = [[col] * col for _ in range(row)]
    topRes = [[-1] * row for _ in range(col)]
    bottomRes = [[row] * row for _ in range(col)]

    for r in range(row):
        for c in range(col):
            if r - 1 >= 0:
                topRes[c][r] = topRes[c][r - 1]
            if c - 1 >= 0:
                leftRes[r][c] = leftRes[r][c - 1]
            if isSpecial[r][c]:
                topRes[c][r] = r
                leftRes[r][c] = c

    for r in range(row - 1, -1, -1):
        for c in range(col - 1, -1, -1):
            if r + 1 < row:
                bottomRes[c][r] = bottomRes[c][r + 1]
            if c + 1 < col:
                rightRes[r][c] = rightRes[r][c + 1]
            if isSpecial[r][c]:
                bottomRes[c][r] = r
                rightRes[r][c] = c

    res = [
        [(topRes[c][r], bottomRes[c][r], leftRes[r][c], rightRes[r][c]) for c in range(col)]
        for r in range(row)
    ]

    return res
