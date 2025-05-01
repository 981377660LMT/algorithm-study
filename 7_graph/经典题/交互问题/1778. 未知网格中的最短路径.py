# https://leetcode.cn/problems/shortest-path-in-a-hidden-grid/
# 你需要找到起点到终点的最短路径，然而你不知道网格的大小、起点和终点。你只能向 GridMaster 对象查询。
# dfs建图 + Dijkstra求成本最小路径


from collections import deque


class GridMaster(object):
    def canMove(self, direction: str) -> bool: ...

    def move(self, direction: str) -> int: ...

    def isTarget(self) -> None: ...


# Directions 的映射：字符 → (dx, dy) 与 相反方向
DIRS = {"U": (-1, 0), "D": (1, 0), "L": (0, -1), "R": (0, 1)}
BACK = {"U": "D", "D": "U", "L": "R", "R": "L"}


class Solution:
    def findShortestPath(self, master: "GridMaster") -> int:
        """
        1) DFS 探索所有可到达格子，建立本地坐标 -> 状态映射
        2) BFS 在这张本地图上求最短路
        """

        # ----------- 1. 探图阶段 -----------
        # 由于题目 |grid| ≤ 5000 而且未知大小，采用 dict 存稀疏坐标
        WALL, EMPTY, TARGET = "#", ".", "T"
        grid = {}  # (x, y) -> 状态; 0,0 为起点
        tgt_coord = None

        def dfs(x: int, y: int) -> None:
            nonlocal tgt_coord
            if master.isTarget():
                grid[(x, y)] = TARGET
                tgt_coord = (x, y)
            else:
                grid[(x, y)] = EMPTY

            for d, (dx, dy) in DIRS.items():
                nx, ny = x + dx, y + dy
                if (nx, ny) in grid:  # 已经探索过
                    continue
                if not master.canMove(d):
                    grid[(nx, ny)] = WALL  # 障碍 / 边界
                    continue

                master.move(d)
                dfs(nx, ny)
                master.move(BACK[d])

        dfs(0, 0)

        if tgt_coord is None:
            return -1

        # ----------- 2. BFS 求最短路 -----------

        q = deque([(0, 0, 0)])  # (x, y, dist)
        seen = {(0, 0)}

        while q:
            x, y, dist = q.popleft()
            if (x, y) == tgt_coord:
                return dist
            for dx, dy in DIRS.values():
                nx, ny = x + dx, y + dy
                if grid.get((nx, ny)) in (EMPTY, TARGET) and (nx, ny) not in seen:
                    seen.add((nx, ny))
                    q.append((nx, ny, dist + 1))

        return -1
