from typing import List, Tuple
from heapq import heappush, heappop

# 由空地（用 0 表示）和墙（用 1 表示）组成的迷宫 maze 中有一个球
# 球可以途经空地向 上、下、左、右 四个方向滚动，且在遇到墙壁前不会停止滚动
# 球可以途经空地向 上、下、左、右 四个方向滚动，且在遇到墙壁前不会停止滚动
# !找出让球停在目的地的最短距离

# 注意:这是带权的bfs(不是所有前进速度一样，不一定先到终点就最短)
# 带权bfs需要dist数组/单源最短路径dijkstra算法

INF = 0x7FFFFFFF


class Solution:
    def shortestDistance(
        self, maze: List[List[int]], start: List[int], destination: List[int]
    ) -> int:
        if not any(maze):
            return -1

        # 启发函数优化
        def valuate(x: int, y: int) -> int:
            tx, ty = destination
            return abs(x - tx) + abs(y - ty)

        m, n = len(maze), len(maze[0])
        queue: List[Tuple[int, int, int]] = [(valuate(0, 0), start[0], start[1])]
        dist = [[INF for _ in range(n)] for _ in range(m)]
        dist[start[0]][start[1]] = 0

        while queue:
            _, r, c = heappop(queue)
            # 判断visited

            for dr, dc in ((0, 1), (1, 0), (0, -1), (-1, 0)):
                nr = r + dr
                nc = c + dc
                steps = 1
                # 沿着这个方向一直进行
                while 0 <= nr < m and 0 <= nc < n and maze[nr][nc] == 0:
                    nr += dr
                    nc += dc
                    steps += 1
                # 碰壁后退
                nr -= dr
                nc -= dc
                steps -= 1
                if dist[r][c] + steps < dist[nr][nc]:
                    dist[nr][nc] = dist[r][c] + steps
                    heappush(queue, (dist[nr][nc] + valuate(nr, nc), nr, nc))

        res = dist[destination[0]][destination[1]]
        return -1 if res == INF else res


print(
    Solution().shortestDistance(
        [[0, 0, 1, 0, 0], [0, 0, 0, 0, 0], [0, 0, 0, 1, 0], [1, 1, 0, 1, 1], [0, 0, 0, 0, 0]],
        [0, 4],
        [4, 4],
    )
)
