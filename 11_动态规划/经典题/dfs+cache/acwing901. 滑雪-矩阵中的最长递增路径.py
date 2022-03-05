# 现在给定你一个二维矩阵表示滑雪场各区域的高度，
# 请你找出在该滑雪场中能够完成的最长下降轨迹，并输出其长度(可经过最大区域数)。
# 1≤R,C≤300,
# 0≤矩阵中整数≤10000
from itertools import product


row, col = map(int, input().split())
mat = [list(map(int, input().split())) for _ in range(row)]


memo = [[-1] * col for _ in range(row)]


def dfs(x: int, y: int) -> int:
    if memo[x][y] != -1:
        return memo[x][y]

    res = 1
    for dx, dy in zip((0, 1, 0, -1), (1, 0, -1, 0)):
        nx, ny = x + dx, y + dy
        if 0 <= nx < row and 0 <= ny < col and mat[nx][ny] < mat[x][y]:
            res = max(res, dfs(nx, ny) + 1)

    memo[x][y] = res
    return res


print(max(dfs(r, c) for r, c in product(range(row), range(col))))
