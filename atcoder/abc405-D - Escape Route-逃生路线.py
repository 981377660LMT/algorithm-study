# abc405-D - Escape Route-逃生路线(多源bfs)
# https://atcoder.jp/contests/abc405/tasks/abc405_d
#
# 在每个通道上写下上、下、左或右方向的箭头，到达最近的逃生出口.

from collections import deque
import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")

DIR4 = [(1, 0), (0, 1), (-1, 0), (0, -1)]
MARK = ["^", "<", "v", ">"]

if __name__ == "__main__":
    H, W = map(int, input().split())
    grid = [list(input()) for _ in range(H)]

    queue = deque()
    for r in range(H):
        for c in range(W):
            if grid[r][c] == "E":
                queue.append((r, c))

    while queue:
        r, c = queue.popleft()
        for i, (dr, dc) in enumerate(DIR4):
            nr, nc = r + dr, c + dc
            if 0 <= nr < H and 0 <= nc < W and grid[nr][nc] == ".":
                grid[nr][nc] = MARK[i]
                queue.append((nr, nc))

    for r in range(H):
        print("".join(grid[r]))
