# C - Sierpinski carpet
# 二维分治(分形地毯)
# https://atcoder.jp/contests/abc357/tasks/abc357_c


from typing import List


def solve(n: int) -> List[List[bool]]:
    dp = [[True]]
    for _ in range(n):
        m = len(dp)
        ndp = [[False] * (3 * m) for _ in range(3 * m)]
        for i in range(3):
            for j in range(3):
                if i == 1 and j == 1:
                    continue
                for x in range(m):
                    for y in range(m):
                        ndp[i * m + x][j * m + y] = dp[x][y]
        dp = ndp

    return dp


if __name__ == "__main__":
    n = int(input())
    grid = solve(n)
    for row in grid:
        print("".join("#" if x else "." for x in row))
