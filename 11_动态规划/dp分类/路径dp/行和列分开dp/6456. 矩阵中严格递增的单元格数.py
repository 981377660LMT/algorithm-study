# 6456. 矩阵中严格递增的单元格数
# https://atcoder.jp/contests/abc224/tasks/abc224_e
# E - Integers on Grid (爬山)
# 有n个方格上有正整数，其余的方格上的数字为0
# 在方格上走，每次只能走同一行或同一列，
# 且到达的格点上的值必须严格大于(真に大きい)当前值，
# 问从给定的n个点出发 最多走多少步
# !ROW,COL<=2e5 n<=2e5


# !1.倒着考虑 求最长路
# !2.因为相等不能取,所以相等的点一起处理
# !3.行和列分开dp rowDp表示当前行的数最多走了多少步 colDp表示当前列的数最多走了多少步
# !dp[(row, col)] = max(rowMax[row], colMax[col]) + 1
# !rowMax[row] = max(rowMax[row], dp[(row, col)])
# !colMax[col] = max(colMax[col], dp[(row, col)])

from collections import defaultdict


INF = int(4e18)

if __name__ == "__main__":
    ROW, COL, n = map(int, input().split())
    mp = defaultdict(list)
    points = []
    for _ in range(n):
        row, col, num = map(int, input().split())
        row, col = row - 1, col - 1
        mp[num].append((row, col))
        points.append((row, col, num))

    keys = sorted(mp, reverse=True)
    dp = defaultdict(int)
    rowMax, colMax = [0] * ROW, [0] * COL
    for key in keys:
        list_ = mp[key]
        for row, col in list_:
            dp[(row, col)] = max(rowMax[row], colMax[col]) + 1
        for row, col in list_:
            rowMax[row] = max(rowMax[row], dp[(row, col)])
            colMax[col] = max(colMax[col], dp[(row, col)])

    for row, col, _ in points:
        print(dp[(row, col)] - 1)
