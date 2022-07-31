from collections import deque
from heapq import heappop, heappush
from typing import List


mapping = {(-1, 0): "^", (1, 0): "v", (0, -1): "<", (0, 1): ">"}
INF = int(1e20)
# LCP 56. 信物传送


class Solution:
    def conveyorBelt(self, matrix: List[str], start: List[int], end: List[int]) -> int:
        grid = [list(row) for row in matrix]
        ROW, COL = len(matrix), len(matrix[0])
        sr, sc = start
        pq = [(0, sr, sc)]

        dist = [[INF] * COL for _ in range(ROW)]
        dist[sr][sc] = 0

        while pq:
            curDist, curRow, curCol = heappop(pq)
            if (curRow, curCol) == (end[0], end[1]):
                return curDist
            # 不变
            for dr, dc in [(0, 1), (0, -1), (1, 0), (-1, 0)]:
                nr, nc = curRow + dr, curCol + dc
                if 0 <= nr < ROW and 0 <= nc < COL:
                    distCand = curDist + int(mapping[(dr, dc)] != grid[curRow][curCol])
                    if dist[nr][nc] > distCand:
                        dist[nr][nc] = distCand
                        heappush(pq, (dist[nr][nc], nr, nc))

        return -1

    def conveyorBelt2(self, matrix: List[str], start: List[int], end: List[int]) -> int:
        """01 bfs 优化"""
        grid = [list(row) for row in matrix]
        ROW, COL = len(matrix), len(matrix[0])
        sr, sc = start
        queue = deque([(0, sr, sc)])

        dist = [[INF] * COL for _ in range(ROW)]
        dist[sr][sc] = 0

        while queue:
            curDist, curRow, curCol = queue.popleft()
            if (curRow, curCol) == (end[0], end[1]):
                return curDist
            # 不变
            for dr, dc in [(0, 1), (0, -1), (1, 0), (-1, 0)]:
                nr, nc = curRow + dr, curCol + dc
                if 0 <= nr < ROW and 0 <= nc < COL:
                    weight = int(mapping[(dr, dc)] != grid[curRow][curCol])
                    distCand = curDist + weight
                    if dist[nr][nc] > distCand:
                        dist[nr][nc] = distCand
                        if weight == 0:
                            queue.appendleft((dist[nr][nc], nr, nc))
                        else:
                            queue.append((dist[nr][nc], nr, nc))
        return -1


print(Solution().conveyorBelt(matrix=[">>v", "v^<", "<><"], start=[0, 1], end=[2, 0]))
