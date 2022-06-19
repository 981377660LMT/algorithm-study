from typing import List, Set, Tuple
from heapq import heappop, heappush

# T:目标 B:箱子 S：起点 #:墙
# 返回将箱子推到目标位置的最小 推动 次数，如果无法做到，请返回 -1。
# 我们只需要返回推箱子的次数。


Position = Tuple[int, int]
DIR4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]


class Solution:
    def minPushBox(self, grid: List[List[str]]) -> int:
        def isInvalid(pos: Position):  # return whether the location is in the grid and not a wall
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

        # 估值距离，箱子推动次数，人，箱子
        pq = [(0, personPos, boxPos)]
        visited: Set[Tuple[Position, Position]] = set()

        while pq:
            moves, person, box = heappop(pq)
            if box == target:
                return moves
            if (person, box) in visited:  # do not visit same state again
                continue
            visited.add((person, box))

            for dr, dc in DIR4:
                nextP = (person[0] + dr, person[1] + dc)
                if isInvalid(nextP):
                    continue
                # !人和箱子重合，就表示推了箱子
                if nextP == box:
                    nextB = (box[0] + dr, box[1] + dc)
                    if isInvalid(nextB):
                        continue
                    heappush(pq, (moves + 1, nextP, nextB))
                else:
                    heappush(pq, (moves, nextP, box))
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
