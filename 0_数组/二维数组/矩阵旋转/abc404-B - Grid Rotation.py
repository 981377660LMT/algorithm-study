import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def rightRotate(matrix: List[List[str]]) -> List[List[str]]:
    """顺时针旋转矩阵90度."""
    return [list(col[::-1]) for col in zip(*matrix)]


def countDiff(grid1: List[List[str]], grid2: List[List[str]]) -> int:
    """计算两个矩阵中不同元素的个数."""
    return sum(1 for row1, row2 in zip(grid1, grid2) for v1, v2 in zip(row1, row2) if v1 != v2)


if __name__ == "__main__":
    N = int(input())
    grid1 = [list(input()) for _ in range(N)]
    grid2 = [list(input()) for _ in range(N)]

    res = INF
    for i in range(4):
        res = min(res, countDiff(grid1, grid2) + i)
        grid1 = rightRotate(grid1)
    print(res)
