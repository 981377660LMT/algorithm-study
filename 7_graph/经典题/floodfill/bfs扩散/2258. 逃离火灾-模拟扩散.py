# 2258. 逃离火灾
# https://leetcode.cn/problems/escape-the-spreading-fire/description/

# 0 表示草地。
# 1 表示着火的格子。
# 2 表示一座墙，你跟火都不能通过这个格子。
# 2 <= m, n <= 300
# 4 <= m * n <= 2e4

# 一开始你在最左上角的格子 (0, 0) ，你想要到达最右下角的安全屋格子 (m - 1, n - 1) 。
# 每一分钟，你可以移动到 相邻 的草地格子。
# 每次你移动 之后 ，着火的格子会扩散到所有不是墙的 相邻 格子。
# !请你返回你在初始位置可以停留的 最多 分钟数，且停留完这段时间后你还能安全到达安全屋。
# 如果无法实现，请你返回 -1 。如果不管你在初始位置停留多久，你 总是 能到达安全屋，请你返回 1e9 。


# 首先通过 BFS 处理出人到每个格子的最短时间 manTime，
# 以及火到每个格子的最短时间 fireTime。


from typing import List
from collections import deque

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]
INF = int(1e9)

# enum
UNVISITED = 0
PERSON_VISITED = 1
FIRE_VISITED = 2
WALL = 3


class Solution:
    def maximumMinutes(self, grid: List[List[int]]) -> int:
        def check(mid: int) -> bool:
            """在初始位置停留mid分钟后,人是否能到达终点"""

            def spread1() -> None:
                """人扩散"""
                len_ = len(queue1)
                for _ in range(len_):
                    curR, curC = queue1.popleft()
                    if visited[curR][curC] != PERSON_VISITED:  # !只能由人访问过的点才能继续扩散
                        continue
                    for dr, dc in DIR4:
                        nextR, nextC = curR + dr, curC + dc
                        if (
                            0 <= nextR < ROW
                            and 0 <= nextC < COL
                            and visited[nextR][nextC] == UNVISITED
                        ):
                            visited[nextR][nextC] = PERSON_VISITED
                            queue1.append((nextR, nextC))

            def spread2() -> None:
                """火扩散"""
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
                        ):
                            visited[nextR][nextC] = FIRE_VISITED
                            queue2.append((nextR, nextC))

            targetR, targetC = ROW - 1, COL - 1
            visited = [[UNVISITED] * COL for _ in range(ROW)]  # 0:未访问,1:人访问过,2:火访问过,3:墙壁
            queue1 = deque()
            queue1.append((0, 0))
            visited[0][0] = PERSON_VISITED
            queue2 = deque()
            for r, row in enumerate(grid):
                for c, cell in enumerate(row):
                    if cell == 1:
                        queue2.append((r, c))
                        visited[r][c] = FIRE_VISITED
                    elif cell == 2:
                        visited[r][c] = WALL

            for _ in range(mid):
                spread2()

            while queue1:
                spread1()  # !如果你到达安全屋后，火马上到了安全屋，这视为你能够安全到达安全屋。
                if visited[targetR][targetC] == PERSON_VISITED:
                    return True
                spread2()

            return False

        ROW, COL = len(grid), len(grid[0])
        left, right = 0, ROW * COL

        if not check(left):  # 无法实现
            return -1
        if check(right):  # 不管你在初始位置停留多久，你 总是 能到达安全屋
            return INF

        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


print(
    Solution().maximumMinutes(
        grid=[
            [0, 2, 0, 0, 0, 0, 0],
            [0, 0, 0, 2, 2, 1, 0],
            [0, 2, 0, 0, 1, 2, 0],
            [0, 0, 2, 2, 2, 0, 2],
            [0, 0, 0, 0, 0, 0, 0],
        ]
    )
)
