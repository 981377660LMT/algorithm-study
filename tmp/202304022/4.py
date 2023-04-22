from functools import lru_cache
from heapq import heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 随着不断的深入，小扣来到了守护者之森寻找的魔法水晶。首先，他必须先通过守护者的考验。

# 考验的区域是一个正方形的迷宫，maze[i][j] 表示在迷宫 i 行 j 列的地形：

# 若为 . ，表示可以到达的空地；
# 若为 # ，表示不可到达的墙壁；
# 若为 S ，表示小扣的初始位置；
# 若为 T ，表示魔法水晶的位置。
# 小扣每次可以向 上、下、左、右 相邻的位置移动一格。而守护者拥有一份「传送魔法卷轴」，使用规则如下：


# 魔法需要在小扣位于 空地 时才能释放，发动后卷轴消失；；
# 发动后，小扣会被传送到水平或者竖直的镜像位置，且目标位置不得为墙壁(如下图所示)；


# 在使用卷轴后，小扣将被「附加负面效果」，因此小扣需要尽可能缩短传送后到达魔法水晶的距离。而守护者的目标是阻止小扣到达魔法水晶的位置；如果无法阻止，则尽可能 增加 小扣传送后到达魔法水晶的距离。
# 假设小扣和守护者都按最优策略行事，返回小扣需要在 「附加负面效果」的情况下 最少 移动多少次才能到达魔法水晶。如果无法到达，返回 -1。

# 4e4


DIR4 = ((-1, 0), (1, 0), (0, -1), (0, 1))


def max(a, b):
    return a if a > b else b


class Solution:
    def challengeOfTheKeeper(self, maze: List[str]) -> int:
        ROW, COL = len(maze), len(maze[0])
        S, T = -1, -1
        for i in range(ROW):
            for j in range(COL):
                if maze[i][j] == "S":
                    S = i * COL + j
                elif maze[i][j] == "T":
                    T = i * COL + j

        def bfs1() -> List[int]:
            """T到各个点的距离"""
            dist = [INF] * (ROW * COL)
            dist[T] = 0
            queue = deque([T])
            while queue:
                cur = queue.popleft()
                x, y = cur // COL, cur % COL
                for dx, dy in DIR4:
                    nx, ny = x + dx, y + dy
                    if 0 <= nx < ROW and 0 <= ny < COL and maze[nx][ny] != "#":
                        if dist[nx * COL + ny] == INF:
                            dist[nx * COL + ny] = dist[cur] + 1
                            queue.append(nx * COL + ny)
            return dist

        distToT = bfs1()

        # 扣的位
        visited = set()
        visited.add(S)
        queue = deque([S])
        while queue:
            cur = queue.popleft()
            x, y = cur // COL, cur % COL
            for dx, dy in DIR4:
                nx, ny = x + dx, y + dy
                if (
                    0 <= nx < ROW
                    and 0 <= ny < COL
                    and maze[nx][ny] != "#"
                    and (nx * COL + ny) not in visited
                ):
                    visited.add(nx * COL + ny)
                    queue.append((nx * COL + ny))

        if T not in visited:
            return -1  # 不连通

        dist2 = [0] * (ROW * COL)  # 到这个点之后，再到T的最大距离
        for i in range(ROW):
            for j in range(COL):
                cur = i * COL + j
                if maze[i][j] == "." and cur in visited:
                    nextRow1, nextCol1 = ROW - 1 - i, j  # 水平镜像
                    if maze[nextRow1][nextCol1] != "#":
                        dist2[cur] = max(dist2[cur], distToT[nextRow1 * COL + nextCol1])
                    nextRow2, nextCol2 = i, COL - 1 - j  # 竖直镜像
                    if maze[nextRow2][nextCol2] != "#":
                        dist2[cur] = max(dist2[cur], distToT[nextRow2 * COL + nextCol2])

        # dijk走一遍到终点
        pq = [(0, S)]
        dist = [INF] * (ROW * COL)
        dist[S] = 0
        while pq:
            d, cur = heappop(pq)
            if d > dist[cur]:
                continue
            if cur == T:
                return d
            x, y = cur // COL, cur % COL
            for dx, dy in DIR4:
                nx, ny = x + dx, y + dy
                nextPos = nx * COL + ny
                if 0 <= nx < ROW and 0 <= ny < COL and maze[nx][ny] != "#":
                    nextCand = max(dist2[nextPos], d)
                    if dist[nextPos] > nextCand:
                        dist[nextPos] = nextCand
                        heappush(pq, (nextCand, nextPos))
        return -1


maze = [".....", "##S..", "...#.", "T.#..", "###.."]
print(Solution().challengeOfTheKeeper(maze))
