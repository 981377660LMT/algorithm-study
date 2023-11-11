"""旋转平移同构判断
两个只含有字符'.'和'#'的矩阵，判断是否可以通过旋转和平移得到同一个矩阵
ROW,COL<=200
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


from typing import List, TypeVar


T = TypeVar("T")


def rotate(matrix: List[List[T]], times: int) -> List[List[T]]:
    """顺时针旋转矩阵90度`times`次"""
    if times == 0:
        return [list(row) for row in matrix]
    res = [list(col[::-1]) for col in zip(*matrix)]
    for _ in range(times - 1):
        res = [list(col[::-1]) for col in zip(*res)]
    return res


def canTranslate(matrix1: List[List[T]], matrix2: List[List[T]]) -> bool:
    """判断matrix1和matrix2是否可以通过平移得到同一个矩阵"""
    ROW, COL = len(matrix1), len(matrix1[0])
    pos1 = [(i, j) for i in range(ROW) for j in range(COL) if matrix1[i][j] == "#"]
    pos2 = [(i, j) for i in range(ROW) for j in range(COL) if matrix2[i][j] == "#"]
    if len(pos1) != len(pos2):
        return False
    offset = set([(x1 - x2, y1 - y2) for (x1, y1), (x2, y2) in zip(pos1, pos2)])
    return len(offset) <= 1


if __name__ == "__main__":
    n = int(input())
    grid1 = [list(input()) for _ in range(n)]
    grid2 = [list(input()) for _ in range(n)]

    for times in range(4):
        rotated = rotate(grid1, times)
        if canTranslate(rotated, grid2):
            print("Yes")
            exit(0)
    print("No")
