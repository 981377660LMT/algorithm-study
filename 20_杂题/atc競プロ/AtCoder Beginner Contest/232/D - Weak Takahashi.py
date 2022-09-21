# 彼が立ち止まるまでに通ることのできるマスは最大で何マスですか
# dp[i][j] = max(dp[i-1][j],dp[i][j-1])+1 if grid[i][j]=='.' else 0
# 1 <= ROW, COL <= 100
# !能通过最多多少个点 bfs/dp
# !二维网格dp还是用dfs 注意dp初始化-INF

from collections import deque
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# '.' '#'
if __name__ == "__main__":

    ROW, COL = map(int, input().split())
    grid = [input() for _ in range(ROW)]

    # # 貰うDP
    # dp = [[-INF] * COL for _ in range(ROW)]
    # dp[0][0] = 1 if grid[0][0] == "." else 0
    # for r in range(ROW):
    #     for c in range(COL):
    #         if grid[r][c] == ".":
    #             if r > 0 and grid[r - 1][c] == ".":
    #                 dp[r][c] = max(dp[r][c], dp[r - 1][c] + 1)
    #             if c > 0 and grid[r][c - 1] == ".":
    #                 dp[r][c] = max(dp[r][c], dp[r][c - 1] + 1)
    # print(max(map(max, dp)))

    # # 配るDP
    # dp = [[-INF] * COL for _ in range(ROW)]
    # dp[0][0] = 1 if grid[0][0] == "." else 0
    # for r in range(ROW):
    #     for c in range(COL):
    #         if grid[r][c] == ".":
    #             if r < ROW - 1 and grid[r + 1][c] == ".":
    #                 dp[r + 1][c] = max(dp[r + 1][c], dp[r][c] + 1)
    #             if c < COL - 1 and grid[r][c + 1] == ".":
    #                 dp[r][c + 1] = max(dp[r][c + 1], dp[r][c] + 1)
    # print(max(map(max, dp)))

    # # dfs
    @lru_cache(None)
    def dfs(r: int, c: int) -> int:
        if (not (0 <= r < ROW and 0 <= c < COL)) or grid[r][c] == "#":
            return 0
        res = 0
        for dr, dc in ((1, 0), (0, 1)):
            nr, nc = r + dr, c + dc
            res = max(res, dfs(nr, nc) + 1)
        return res

    print(dfs(0, 0))

    # bfs 求最长路
    # queue = deque([(0, 0, 0)])
    # dist = [[-INF] * COL for _ in range(ROW)]
    # dist[0][0] = 0
    # while queue:
    #     curRow, curCol, curDist = queue.popleft()
    #     if dist[curRow][curCol] < curDist:
    #         continue
    #     for dr, dc in ((1, 0), (0, 1)):
    #         nr, nc = curRow + dr, curCol + dc
    #         if (not (0 <= nr < ROW and 0 <= nc < COL)) or grid[nr][nc] == "#":
    #             continue
    #         if dist[nr][nc] < curDist + 1:
    #             dist[nr][nc] = curDist + 1
    #             queue.append((nr, nc, curDist + 1))
    # print(max(map(max, dist)) + 1)
