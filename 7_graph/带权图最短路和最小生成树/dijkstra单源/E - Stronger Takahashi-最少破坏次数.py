# '.':路 '#':障碍物
# 每次可以破坏周围任意的2*2的区域内的墙壁，变成路
# !求(0,0)走到(n-1,m-1)的最少破坏次数
# ROW,COL<=500

# !注意: 破坏等价于使得2*2的格子和当前格子邻接
from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]

if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [input() for _ in range(ROW)]

    dist = [[INF] * COL for _ in range(ROW)]
    dist[0][0] = 0

    pq = [(0, 0, 0)]
    while pq:
        curDist, curRow, curCol = heappop(pq)
        if curDist > dist[curRow][curCol]:
            continue
        if curRow == ROW - 1 and curCol == COL - 1:
            print(curDist)
            exit(0)

        # !不出拳
        for dr, dc in DIR4:
            nextRow, nextCol = curRow + dr, curCol + dc
            if (0 <= nextRow < ROW and 0 <= nextCol < COL) and (grid[nextRow][nextCol] == "."):
                if curDist < dist[nextRow][nextCol]:
                    dist[nextRow][nextCol] = curDist
                    heappush(pq, (curDist, nextRow, nextCol))

        # !出拳后可以到达的位置
        for dr in range(-2, 3):
            for dc in range(-2, 3):
                if abs(dr) + abs(dc) > 3:
                    continue
                nextRow, nextCol = curRow + dr, curCol + dc
                if 0 <= nextRow < ROW and 0 <= nextCol < COL:
                    cand = curDist + 1
                    if cand < dist[nextRow][nextCol]:
                        dist[nextRow][nextCol] = cand
                        heappush(pq, (cand, nextRow, nextCol))
