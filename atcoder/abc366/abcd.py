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


if __name__ == "__main__":
    N = int(input())
    S = [input() for _ in range(N)]
    M = max(len(s) for s in S)
    grid = [["*"] * M for _ in range(N)]
    for i, s in enumerate(S):
        for j, c in enumerate(s):
            grid[i][j] = c
    trans = rotate(grid, 1)
    for row in trans:
        cur = "".join(row).rstrip("*")
        print(cur)
