# 传送带 问最后停在哪个位置 如果循环 输出-1
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

DIR = {"U": (-1, 0), "D": (1, 0), "L": (0, -1), "R": (0, 1)}
ROW, COL = map(int, input().split())

grid = []
for _ in range(ROW):
    row = list(input())
    grid.append(row)

curRow, curCol = 0, 0
visited = [[False] * COL for _ in range(ROW)]

while True:
    visited[curRow][curCol] = True
    dr, dc = DIR[grid[curRow][curCol]]
    nextRow, nextCol = curRow + dr, curCol + dc
    if nextRow < 0 or nextRow >= ROW or nextCol < 0 or nextCol >= COL:
        break
    if visited[nextRow][nextCol]:
        print(-1)
        exit(0)
    curRow, curCol = nextRow, nextCol

print(curRow + 1, curCol + 1)
