import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    h, w, d = map(int, input().split())
    grid = [input() for _ in range(h)]

    from collections import deque

    visited = [[-1] * w for _ in range(h)]
    queue = deque()

    for i in range(h):
        for j in range(w):
            if grid[i][j] == "H":
                visited[i][j] = 0
                queue.append((i, j))

    directions = [(-1, 0), (1, 0), (0, -1), (0, 1)]

    while queue:
        x, y = queue.popleft()
        if visited[x][y] >= d:
            continue
        for dx, dy in directions:
            nx, ny = x + dx, y + dy
            if 0 <= nx < h and 0 <= ny < w:
                if grid[nx][ny] != "#" and visited[nx][ny] == -1:
                    visited[nx][ny] = visited[x][y] + 1
                    queue.append((nx, ny))

    result = 0
    for i in range(h):
        for j in range(w):
            if (
                (grid[i][j] == "." or grid[i][j] == "H")
                and visited[i][j] != -1
                and visited[i][j] <= d
            ):
                result += 1

    print(result)
