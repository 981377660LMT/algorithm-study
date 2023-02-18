# https://maspypy.github.io/library/other/magic_square.hpp
# magicSquare
# n*n 的幻方(河图) 行/列/对角线之和相等


from typing import List


def 河图(n: int) -> List[List[int]]:
    assert n % 2 == 1
    res = [[0] * n for _ in range(n)]
    x, y = 0, n // 2
    for i in range(n * n):
        res[x][y] = i + 1
        nx = n - 1 if x == 0 else x - 1
        ny = 0 if y == n - 1 else y + 1
        if res[nx][ny] != 0:
            x = 0 if x == n - 1 else x + 1
            nx, ny = x, y
        x, y = nx, ny
    return res


print(河图(3))
print(河图(5))
