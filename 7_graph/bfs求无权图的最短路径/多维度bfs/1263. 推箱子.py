from collections import deque
from typing import Deque, List, Set, Tuple


# T:目标 B:箱子 S：起点 #:墙
# !返回将箱子推到目标位置的最小 推动 次数，如果无法做到，请返回 -1。

# 1.怎么保证不重复访问:两个位置作为状态
# 2.怎么处理推箱子的位置和不推箱子位置的先后顺序：01bfs

Position = Tuple[int, int]
DIR4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]


class Solution:
    def minPushBox(self, grid: List[List[str]]) -> int:
        """01bfs求两种权值的最短路 时空复杂度O(n^2*m^2)"""

        def isInvalid(pos: Position) -> bool:
            r, c = pos
            if r < 0 or r >= ROW:
                return True
            if c < 0 or c >= COL:
                return True
            return grid[r][c] == "#"

        ROW, COL = len(grid), len(grid[0])
        target, boxPos, personPos = (0, 0), (0, 0), (0, 0)
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == "T":
                    target = (r, c)
                if grid[r][c] == "B":
                    boxPos = (r, c)
                if grid[r][c] == "S":
                    personPos = (r, c)

        # (可以加一个估值距离)，箱子推动次数，人，箱子
        queue: Deque[Tuple[int, Position, Position]] = deque([(0, personPos, boxPos)])
        visited: Set[Tuple[Position, Position]] = set()

        while queue:
            boxMove, personPos, boxPos = queue.popleft()
            if boxPos == target:
                return boxMove
            if (personPos, boxPos) in visited:
                continue
            visited.add((personPos, boxPos))

            for dr, dc in DIR4:
                nextPerson = (personPos[0] + dr, personPos[1] + dc)
                if isInvalid(nextPerson):
                    continue
                # !人和箱子重合，就表示推了箱子
                if nextPerson == boxPos:
                    nextB = (boxPos[0] + dr, boxPos[1] + dc)
                    if isInvalid(nextB):
                        continue
                    queue.append((boxMove + 1, nextPerson, nextB))
                else:
                    queue.appendleft((boxMove, nextPerson, boxPos))

        return -1


print(
    Solution().minPushBox(
        [
            ["#", "#", "#", "#", "#", "#"],
            ["#", "T", "#", "#", "#", "#"],
            ["#", ".", ".", "B", ".", "#"],
            ["#", ".", "#", "#", ".", "#"],
            ["#", ".", ".", ".", "S", "#"],
            ["#", "#", "#", "#", "#", "#"],
        ]
    )
)
