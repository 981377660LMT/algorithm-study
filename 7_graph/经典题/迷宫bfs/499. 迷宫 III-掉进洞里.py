from typing import Deque, List, Tuple
from collections import deque

# 由空地（用 0 表示）和墙（用 1 表示）组成的迷宫 maze 中有一个球
# 球可以途经空地向 上、下、左、右 四个方向滚动，且在遇到墙壁前不会停止滚动
# 球可以途经空地向 上、下、左、右 四个方向滚动，且在遇到墙壁前不会停止滚动
# 找出让球以最短距离掉进洞里的路径
# 可能有多条最短路径， 请输出字典序最小的路径。


# 注意:这是带权的bfs(不是所有前进速度一样，不一定先到终点就最短)
# 带权bfs需要dist数组/单源最短路径dijkstra算法
# 开一个数组二维记录每个点处的path


INF = 0x7FFFFFFF


class Solution:
    def findShortestWay(self, maze: List[List[int]], ball: List[int], hole: List[int]) -> str:
        if not any(maze):
            return "impossible"

        m, n = len(maze), len(maze[0])
        queue: Deque[Tuple[int, int]] = deque()
        dist = [[INF for _ in range(n)] for _ in range(m)]
        path = [["impossible" for _ in range(n)] for _ in range(m)]

        dist[ball[0]][ball[1]] = 0
        path[ball[0]][ball[1]] = ''
        queue.append((ball[0], ball[1]))

        while queue:
            r, c = queue.popleft()
            for dr, dc, direction in ((0, 1, 'r'), (1, 0, 'd'), (0, -1, 'l'), (-1, 0, 'u')):
                nr = r + dr
                nc = c + dc
                all_step = dist[r][c] + 1
                new_path = path[r][c] + direction

                # 沿着这个方向一直进行
                while (
                    0 <= nr < m
                    and 0 <= nc < n
                    and maze[nr][nc] == 0
                    and not (nr - dr == hole[0] and nc - dc == hole[1])  # 不能是从洞滚过来的(上一步就已经到洞里了)
                ):
                    nr += dr
                    nc += dc
                    all_step += 1

                # 碰壁后退
                nr -= dr
                nc -= dc
                all_step -= 1

                # 更短，或者字典序更小
                if all_step < dist[nr][nc] or all_step == dist[nr][nc] and new_path < path[nr][nc]:
                    dist[nr][nc] = all_step
                    path[nr][nc] = new_path
                    if not [nr, nc] == hole:
                        queue.append((nr, nc))

        return path[hole[0]][hole[1]]


print(
    Solution().findShortestWay(
        [[0, 0, 0, 0, 0], [1, 1, 0, 0, 1], [0, 0, 0, 0, 0], [0, 1, 0, 0, 1], [0, 1, 0, 0, 0]],
        [4, 3],
        [0, 1],
    )
)

