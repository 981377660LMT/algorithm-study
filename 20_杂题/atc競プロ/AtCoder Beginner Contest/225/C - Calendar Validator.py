# 日历检验
# !给定一个行数无穷(1e100)列数为 7 的矩阵，判断给定矩阵是否是这个大矩阵的一部分
# !大矩阵(i,j)处元素为 (i-1)*7+j (1<=i<=1e100, 1<=j<=7)

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [tuple(map(int, input().split())) for _ in range(ROW)]

    # 行的差是否合法
    for row in grid:
        for pre, cur in zip(row, row[1:]):
            if cur - pre != 1:
                print("No")
                exit(0)

    # 列的差是否合法
    for col in zip(*grid):
        for pre, cur in zip(col, col[1:]):
            if cur - pre != 7:
                print("No")
                exit(0)

    # !元素范围是否合法
    sr, sc = divmod(grid[0][0] - 1, 7)
    if sc + COL > 7:
        print("No")
        exit(0)

    print("Yes")
