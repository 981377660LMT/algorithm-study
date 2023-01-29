# 5
# 8 8 8 7 7
# 7 7 8 8 7
# 7 7 7 7 7
# 7 8 8 7 8
# 7 8 8 8 8

# 2 1
from collections import deque
from typing import Tuple


n = int(input())
mat = [list(map(int, input().split())) for _ in range(n)]
visited = [[0] * n for _ in range(n)]

dxs, dys = [-1, -1, 0, 1, 1, 1, 0, -1], [0, -1, -1, -1, 0, 1, 1, 1]

# 1≤n≤1000


def bfs(x: int, y: int) -> Tuple[bool, bool]:
    """bfs floodfill看所有高度相等的点"""
    has_higher = has_lower = False
    queue = deque([(x, y)])
    visited[x][y] = True
    while queue:
        x, y = queue.popleft()
        cv = mat[x][y]
        for dx, dy in zip(dxs, dys):
            nx, ny = x + dx, y + dy
            if 0 <= nx < n and 0 <= ny < n:
                nv = mat[nx][ny]
                if nv != cv:
                    if nv > cv:
                        has_higher = True
                    else:
                        has_lower = True
                elif not visited[nx][ny]:
                    queue.append((nx, ny))
                    visited[nx][ny] = True

    return has_higher, has_lower


def count():
    peak = valley = 0

    for i in range(n):
        for j in range(n):
            if not visited[i][j]:
                has_higher, has_lower = bfs(i, j)
                if not has_higher:
                    peak += 1
                if not has_lower:
                    valley += 1

    print(peak, valley)


count()

