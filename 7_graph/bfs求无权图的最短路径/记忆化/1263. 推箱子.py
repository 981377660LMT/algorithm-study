from typing import List, Set, Tuple
from heapq import heappop, heappush

# T:目标 B:箱子 S：起点 #:墙
# 返回将箱子推到目标位置的最小 推动 次数，如果无法做到，请返回 -1。
# 我们只需要返回推箱子的次数。


Position = Tuple[int, int]


class Solution:
    def minPushBox(self, grid: List[List[str]]) -> int:
        row, col = len(grid), len(grid[0])
        for r in range(row):
            for c in range(col):
                if grid[r][c] == "T":
                    target = (r, c)
                if grid[r][c] == "B":
                    start_box = (r, c)
                if grid[r][c] == "S":
                    start_person = (r, c)

        # 启发估值函数
        def heuristic(box: Position):
            return abs(target[0] - box[0]) + abs(target[1] - box[1])

        def out_bounds(pos: Position):  # return whether the location is in the grid and not a wall
            r, c = pos
            if r < 0 or r >= row:
                return True
            if c < 0 or c >= col:
                return True
            return grid[r][c] == "#"

        # 估值距离，箱子推动次数，人，箱子
        pq = [(heuristic(start_box) + 0, 0, start_person, start_box)]
        visited: Set[Tuple[Position, Position]] = set()

        while pq:
            _, moves, person, box = heappop(pq)
            if box == target:
                return moves
            if (person, box) in visited:  # do not visit same state again
                continue
            visited.add((person, box))

            for dr, dc in [[0, 1], [1, 0], [-1, 0], [0, -1]]:
                new_person = (person[0] + dr, person[1] + dc)
                if out_bounds(new_person):
                    continue
                if new_person == box:
                    new_box = (box[0] + dr, box[1] + dc)
                    if out_bounds(new_box):
                        continue
                    heappush(pq, (heuristic(new_box) + moves + 1, moves + 1, new_person, new_box))
                else:
                    heappush(pq, (heuristic(box) + moves, moves, new_person, box))
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
