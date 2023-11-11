"""
从左上角出发,每次只能向右或向下走一步,最后到达右下角
移动的费用为`走过的格子中最大的k个数的和`
求移动的最小费用
ROW,COL<=30 nums[i][j]<=1e9

1. 按照常规的dp思路,发现需要记录所有经过的格子,具有后效性
2. 解决 dp 的后效性,直接`枚举前k大中的最小值`,这样每次经过一个格子,只需要判断要不要加上这个格子的值
   dp[row][col][min_]
   !枚举某个参数以消除dp的后效性
"""
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    # !dfs O(ROW^2*COL^2) 优化:进行了ROW*COL次O(ROW*COL) 与k无关
    # ROW, COL, k = map(int, input().split())
    # grid = [list(map(int, input().split())) for _ in range(ROW)]

    # @lru_cache(None)
    # def dfs(row: int, col: int, min_: int) -> int:
    #     if row == ROW - 1 and col == COL - 1:
    #         return 0
    #     res = INF
    #     if row + 1 < ROW:
    #         res = min(res, dfs(row + 1, col, min_) + max(0, grid[row + 1][col] - min_))
    #     if col + 1 < COL:
    #         res = min(res, dfs(row, col + 1, min_) + max(0, grid[row][col + 1] - min_))
    #     return res

    # res = INF
    # # 枚举前k大中的最小值
    # for min_ in set(grid[row][col] for row in range(ROW) for col in range(COL)):
    #     res = min(res, dfs(0, 0, min_) + max(0, grid[0][0] - min_) + k * min_)
    # print(res)
    ############################################
    # !dp O(ROW^2*COL^2*k)
    # !dp[row][col][count] 表示从左上角到(row,col)的最小费用,有count个数大于等于min_
    ROW, COL, k = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(ROW)]

    res = INF

    # !枚举前k大中的最小值
    for min_ in set(grid[row][col] for row in range(ROW) for col in range(COL)):
        dp = [[[INF] * (k + 1) for _ in range(COL)] for _ in range(ROW)]
        if grid[0][0] >= min_:
            dp[0][0][1] = grid[0][0]
        if grid[0][0] <= min_:
            dp[0][0][0] = 0
        for r in range(ROW):
            for c in range(COL):
                for count in range(k + 1):
                    if r + 1 < ROW:
                        if count != k and grid[r + 1][c] >= min_:  # !加上这个格子
                            dp[r + 1][c][count + 1] = min(
                                dp[r + 1][c][count + 1], dp[r][c][count] + grid[r + 1][c]
                            )
                        if grid[r + 1][c] <= min_:  # !不加上这个格子
                            dp[r + 1][c][count] = min(dp[r + 1][c][count], dp[r][c][count])
                    if c + 1 < COL:
                        if count != k and grid[r][c + 1] >= min_:
                            dp[r][c + 1][count + 1] = min(
                                dp[r][c + 1][count + 1], dp[r][c][count] + grid[r][c + 1]
                            )
                        if grid[r][c + 1] <= min_:
                            dp[r][c + 1][count] = min(dp[r][c + 1][count], dp[r][c][count])

        res = min(res, dp[ROW - 1][COL - 1][k])

    print(res)
