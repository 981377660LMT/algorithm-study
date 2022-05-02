from typing import List, Tuple
from collections import defaultdict, deque


MOD = int(1e9 + 7)
INF = int(1e20)

# 0 表示草地。
# 1 表示着火的格子。
# 2 表示一座墙，你跟火都不能通过这个格子。
# 2 <= m, n <= 300
# 4 <= m * n <= 2 * 1e4

dirs = [(0, 1), (1, 0), (0, -1), (-1, 0)]


class Solution:
    def maximumMinutes(self, grid: List[List[int]]) -> int:
        def bfs(bad: set) -> bool:
            """是否可以到达"""
            q = deque([(0, 0)])
            visitedQ = set([(0, 0)])
            fireQ = deque(bad)
            while q:
                len1 = len(q)
                for _ in range(len1):
                    r, c = q.popleft()
                    if r == ROW - 1 and c == COL - 1:
                        if (r, c) in bad:
                            self.res = True
                        return True
                    for dr, dc in dirs:
                        nr, nc = r + dr, c + dc
                        if (
                            0 <= nr < ROW
                            and 0 <= nc < COL
                            and grid[nr][nc] == 0
                            and (nr, nc) not in bad
                            and (nr, nc) not in visitedQ
                        ):
                            q.append((nr, nc))
                            visitedQ.add((nr, nc))
                len2 = len(fireQ)
                for _ in range(len2):
                    r, c = fireQ.popleft()
                    for dr, dc in dirs:
                        nr, nc = r + dr, c + dc
                        if (
                            0 <= nr < ROW
                            and 0 <= nc < COL
                            and grid[nr][nc] == 0
                            and (nr, nc) not in bad
                        ):
                            fireQ.append((nr, nc))
                            bad.add((nr, nc))
            return False

        def check(mid: int) -> bool:
            """预处理火的开始位置"""
            curBad = fire.copy()
            queue = deque(curBad)
            for _ in range(mid):
                len_ = len(queue)
                for _ in range(len_):
                    r, c = queue.popleft()
                    for dr, dc in dirs:
                        nr, nc = r + dr, c + dc
                        if (
                            0 <= nr < ROW
                            and 0 <= nc < COL
                            and grid[nr][nc] == 0
                            and (nr, nc) not in curBad
                        ):
                            queue.append((nr, nc))
                            curBad.add((nr, nc))
            return bfs(curBad)

        ROW, COL = len(grid), len(grid[0])
        fire = set()
        for i in range(ROW):
            for j in range(COL):
                if grid[i][j] == 1:
                    fire.add((i, j))

        self.res = False
        left, right = 0, ROW + COL
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                self.res = True
                left = mid + 1
            else:
                right = mid - 1
        if not self.res:
            return -1
        if right >= ROW + COL:
            return int(1e9)

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
print(Solution().maximumMinutes(grid=[[0, 0, 0, 0], [0, 1, 2, 0], [0, 2, 0, 0]]))

# 0
print(
    Solution().maximumMinutes(
        grid=[[0, 2, 0, 0, 1], [0, 2, 0, 2, 2], [0, 2, 0, 0, 0], [0, 0, 2, 2, 0], [0, 0, 0, 0, 0]]
    )
)
