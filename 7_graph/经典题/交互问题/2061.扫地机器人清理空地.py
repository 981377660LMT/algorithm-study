from typing import List

# 起初它位于 [0, 0] 位置，面朝右侧，0 为空地，1 为障碍物，只有当它当前朝向遇到了障碍物或者墙壁时，才可以将朝向顺时针旋转 90 度，
# 问你它最多能清理多少空地。


D = [(0, 1), (1, 0), (0, -1), (-1, 0)]


class Solution:
    def numberOfCleanRooms(self, room: List[List[int]]) -> int:
        m, n = len(room), len(room[0])
        x, y, dr, visited = 0, 0, 0, set()

        while True:
            if (x, y, dr) in visited:
                break
            visited.add((x, y, dr))

            if (
                0 <= x + D[dr][0] < m
                and 0 <= y + D[dr][1] < n
                and room[x + D[dr][0]][y + D[dr][1]] == 0
            ):
                x, y = x + D[dr][0], y + D[dr][1]
            else:
                dr = (dr + 1) % 4

        return len(set((x, y) for x, y, _ in visited))

