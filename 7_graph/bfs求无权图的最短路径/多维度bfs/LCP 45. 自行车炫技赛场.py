# 选手初始位于坐标 position 处且初始速度为 1，请问选手可以`刚好`到其他哪些位置时速度依旧为 1。
# 自行车从高度为 h1 且减速值为 o1 的位置到高度为 h2 且减速值为 o2 的相邻位置（上下左右四个方向），速度变化值为 h1-h2-o2（负值减速，正值增速）。

from typing import List
from collections import deque


class Solution:
    def bicycleYard(
        self, position: List[int], terrain: List[List[int]], obstacle: List[List[int]]
    ) -> List[List[int]]:
        ROW, COL = len(terrain), len(terrain[0])
        sr, sc = position[0], position[1]

        queue = deque([(sr, sc, 1)])
        visited = [[set() for _ in range(COL)] for _ in range(ROW)]
        res = set()
        while queue:
            r, c, curSpeed = queue.popleft()
            if curSpeed in visited[r][c]:
                continue
            visited[r][c].add(curSpeed)
            if (curSpeed == 1) and (r, c) != (position[0], position[1]):
                res.add((r, c))
            for nr, nc in [(r - 1, c), (r + 1, c), (r, c - 1), (r, c + 1)]:
                if 0 <= nr < ROW and 0 <= nc < COL:
                    cand = terrain[r][c] - terrain[nr][nc] - obstacle[nr][nc] + curSpeed
                    if cand > 0:
                        queue.append((nr, nc, cand))  # type: ignore

        return sorted(res)


# 输出：[[0,1]]
# 解释：
# 选手从 [1,1] 处的位置出发，到 [0,1] 处的位置时恰好速度为 1。
print(Solution().bicycleYard(position=[1, 1], terrain=[[5, 0], [0, 6]], obstacle=[[0, 6], [7, 0]]))
