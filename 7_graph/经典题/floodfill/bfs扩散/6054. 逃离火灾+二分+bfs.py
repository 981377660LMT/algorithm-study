from typing import Deque, List, Set, Tuple
from collections import deque


# 0 表示草地。
# 1 表示着火的格子。
# 2 表示一座墙，你跟火都不能通过这个格子。
# 2 <= m, n <= 300
# 4 <= m * n <= 2 * 1e4

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]


class Solution:
    def maximumMinutes(self, grid: List[List[int]]) -> int:
        """时间复杂度O(mnlog(mn))"""

        def bfs(fireVisited: Set[Tuple[int, int]], fireQueue: Deque[Tuple[int, int]]) -> bool:
            """人先走，是否可以到达"""
            queue = deque([(0, 0)])
            visited = set([(0, 0)])

            while queue:
                len1 = len(queue)
                for _ in range(len1):
                    r, c = queue.popleft()
                    if r == ROW - 1 and c == COL - 1:
                        return True

                    # 注意这里啊
                    if (r, c) in fireVisited:
                        continue

                    for dr, dc in DIR4:
                        nr, nc = r + dr, c + dc
                        if (
                            0 <= nr < ROW
                            and 0 <= nc < COL
                            and grid[nr][nc] == 0
                            and (nr, nc) not in fireVisited
                            and (nr, nc) not in visited
                        ):
                            queue.append((nr, nc))
                            visited.add((nr, nc))

                len2 = len(fireQueue)
                for _ in range(len2):
                    r, c = fireQueue.popleft()
                    for dr, dc in DIR4:
                        nr, nc = r + dr, c + dc
                        if (
                            0 <= nr < ROW
                            and 0 <= nc < COL
                            and grid[nr][nc] == 0
                            and (nr, nc) not in fireVisited
                        ):
                            fireQueue.append((nr, nc))
                            fireVisited.add((nr, nc))
            return False

        def check(mid: int) -> bool:
            """预处理火的开始位置,处理出visited和queue"""
            fireVisited = fire.copy()
            fireQueue = deque(fireVisited)

            for _ in range(mid):
                len_ = len(fireQueue)
                for _ in range(len_):
                    r, c = fireQueue.popleft()
                    for dr, dc in DIR4:
                        nr, nc = r + dr, c + dc
                        if (
                            0 <= nr < ROW
                            and 0 <= nc < COL
                            and grid[nr][nc] != 2
                            and (nr, nc) not in fireVisited
                        ):
                            fireQueue.append((nr, nc))
                            fireVisited.add((nr, nc))

            return bfs(fireVisited, fireQueue)

        ROW, COL = len(grid), len(grid[0])
        fire = set((r, c) for r in range(ROW) for c in range(COL) if grid[r][c] == 1)

        if not check(0):
            return -1

        upper = ROW * COL
        left, right = 0, upper
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1

        if right >= upper:
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

# -1
print(Solution().maximumMinutes(grid=[[0, 0, 0, 0], [0, 1, 2, 0], [0, 2, 0, 0]]))

# 0
print(
    Solution().maximumMinutes(
        grid=[[0, 2, 0, 0, 1], [0, 2, 0, 2, 2], [0, 2, 0, 0, 0], [0, 0, 2, 2, 0], [0, 0, 0, 0, 0]]
    )
)
