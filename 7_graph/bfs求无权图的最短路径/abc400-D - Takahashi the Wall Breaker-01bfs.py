# abc400-D - Takahashi the Wall Breaker - 01bfs
# https://atcoder.jp/contests/abc400/tasks/abc400_d
#
#
# 1-移动到四个方向相邻的区域，并且该区域是道路。
# 2-前踢，将前面两个区域的墙变为道路。
# 求(A,B)到(C,D)的前踢次数最小值.


from collections import deque


INF = int(1e18)
DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]

if __name__ == "__main__":
    H, W = map(int, input().split())
    grid = [input() for _ in range(H)]
    A, B, C, D = map(int, input().split())
    A, B, C, D = A - 1, B - 1, C - 1, D - 1

    dist = [[INF] * W for _ in range(H)]
    queue = deque()

    def isIn(x: int, y: int) -> bool:
        return 0 <= x < H and 0 <= y < W

    def update(nx: int, ny: int, nd: int, cost: int) -> None:
        if dist[nx][ny] > nd:
            dist[nx][ny] = nd
            if cost == 0:
                queue.appendleft((nx, ny))
            else:
                queue.append((nx, ny))

    update(A, B, 0, 0)

    while queue:
        x, y = queue.popleft()
        for dx, dy in DIR4:
            nx, ny = x + dx, y + dy
            if isIn(nx, ny):
                if grid[nx][ny] == "#":
                    update(nx, ny, dist[x][y] + 1, 1)
                else:
                    update(nx, ny, dist[x][y], 0)
            nx, ny = x + dx * 2, y + dy * 2
            if isIn(nx, ny) and grid[nx][ny] == "#":
                update(nx, ny, dist[x][y] + 1, 1)

    res = dist[C][D]
    print(res)
