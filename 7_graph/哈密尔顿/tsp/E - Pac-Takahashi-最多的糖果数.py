# 二维迷宫，有一个起点和一个终点，有墙，有最多18个糖果。
# 问从起点到终点，移动距离不超过 t的情况下，能获得的糖果数量的最大值。
# 如果不能到达终点，输出 -1。
# "."=>空地
# "#"=>墙
# "S"=>起点
# "G"=>终点
# "o"=>糖果

# 1<=ROW,COL<=300
# 1<=T<=2e6

# !dp[n][2^n] 表示状态:当前看到第i个糖果，已经拿到的糖果集合.
# 无论访问的顺序如何，只要访问的点集是一样的，那么糖果数就是一样的。
# !通过dp将O(n!)降低到O(2^n*n^2)


from collections import deque
from typing import List

INF = int(1e18)
DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


def pacTakahashi(grid: List[str], limit: int) -> int:
    def bfs(sr: int, sc: int) -> List[List[int]]:
        dist = [[INF] * COL for _ in range(ROW)]
        dist[sr][sc] = 0
        queue = deque([(sr, sc)])
        while queue:
            curRow, curCol = queue.popleft()
            for dr, dc in DIR4:
                nr, nc = curRow + dr, curCol + dc
                if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] != "#":
                    cand = dist[curRow][curCol] + 1
                    if dist[nr][nc] > cand:
                        dist[nr][nc] = cand
                        queue.append((nr, nc))
        return dist

    ROW, COL = len(grid), len(grid[0])
    sr, sc = -1, -1
    tr, tc = -1, -1
    candies = []
    for i in range(ROW):
        for j in range(COL):
            if grid[i][j] == "o":
                candies.append((i, j))
            elif grid[i][j] == "S":
                sr, sc = i, j
            elif grid[i][j] == "G":
                tr, tc = i, j

    if bfs(sr, sc)[tr][tc] > limit:
        print(-1)
        exit(0)

    dists = [bfs(r, c) for r, c in candies]
    C = len(candies)
    dp = [[INF] * (1 << C) for _ in range(C)]
    for i in range(C):
        dp[i][1 << i] = dists[i][sr][sc]
    for state in range(1, 1 << C):
        for last in range(C):
            for next, (r2, c2) in enumerate(candies):
                if not state & (1 << next):
                    nextState = state | (1 << next)
                    cand = dp[last][state] + dists[last][r2][c2]
                    if dp[next][nextState] > cand:
                        dp[next][nextState] = cand

    res = 0
    for state in range(1, 1 << C):
        for last in range(C):
            d = dp[last][state] + dists[last][tr][tc]
            if d <= limit:
                res = max(res, bin(state).count("1"))
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    ROW, COL, t = map(int, input().split())
    grid = [input() for _ in range(ROW)]
    res = pacTakahashi(grid, t)
    print(res)
