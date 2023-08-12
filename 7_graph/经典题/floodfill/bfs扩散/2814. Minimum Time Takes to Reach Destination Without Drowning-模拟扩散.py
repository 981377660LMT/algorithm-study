# https://leetcode.cn/problems/minimum-time-takes-to-reach-destination-without-drowning/

# "S": "起点",
# "D": "终点",
# "*": "水",
# ".": "陆地",
# "X": "墙",
# !注意终点不会被水淹没


# !模拟洪水扩散

from collections import deque
from typing import List

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]
INF = int(1e18)

# enum
UNVISITED = 0
PERSON_VISITED = 1
FIRE_VISITED = 2
WALL = 3


class Solution:
    def minimumSeconds(self, land: List[List[str]]) -> int:
        def spread1() -> None:
            """人扩散"""
            len_ = len(queue1)
            for _ in range(len_):
                curR, curC = queue1.popleft()
                if visited[curR][curC] != PERSON_VISITED:  # !只能由人访问过的点才能继续扩散
                    continue
                for dr, dc in DIR4:
                    nextR, nextC = curR + dr, curC + dc
                    if 0 <= nextR < ROW and 0 <= nextC < COL and visited[nextR][nextC] == UNVISITED:
                        visited[nextR][nextC] = PERSON_VISITED
                        queue1.append((nextR, nextC))

        def spread2() -> None:
            """洪水扩散,洪水无法淹没终点"""
            len_ = len(queue2)
            for _ in range(len_):
                curR, curC = queue2.popleft()
                for dr, dc in DIR4:
                    nextR, nextC = curR + dr, curC + dc
                    if (
                        0 <= nextR < ROW
                        and 0 <= nextC < COL
                        and (
                            visited[nextR][nextC] in (UNVISITED, PERSON_VISITED)
                        )  # !只能继续扩散到未访问/人访问过的点
                        and (nextR, nextC) != (targetR, targetC)  # !终点不会被水淹没
                    ):
                        visited[nextR][nextC] = FIRE_VISITED
                        queue2.append((nextR, nextC))

        ROW, COL = len(land), len(land[0])
        targetR, targetC = -1, -1
        visited = [[UNVISITED] * COL for _ in range(ROW)]  # 0:未访问,1:人访问过,2:洪水访问过,3:墙壁
        queue1 = deque()
        queue2 = deque()
        for r, row in enumerate(land):
            for c, cell in enumerate(row):
                if cell == ".":
                    continue
                elif cell == "*":
                    queue2.append((r, c))
                    visited[r][c] = FIRE_VISITED
                elif cell == "X":
                    visited[r][c] = WALL
                elif cell == "S":
                    queue1.append((r, c))
                    visited[r][c] = PERSON_VISITED
                elif cell == "D":
                    targetR, targetC = r, c

        res = 1
        while queue1:
            spread1()
            if visited[targetR][targetC] == PERSON_VISITED:
                return res
            spread2()
            res += 1

        return -1


print(Solution().minimumSeconds(land=[["D", ".", "*"], [".", ".", "."], [".", "S", "."]]))
