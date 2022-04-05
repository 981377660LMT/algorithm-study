# 每次你把与(0,0)相连的相同颜色的点都置0或1
# 求所有格子变为同色的最少操作数

# 这个过程很像bfs的波纹扩张
from collections import deque


class Solution:
    def solve(self, matrix):
        row, col = len(matrix), len(matrix[0])

        res = 0
        visited = set()
        queue = deque([(0, 0)])
        color = matrix[0][0]

        while len(visited) < row * col:
            res += 1
            nextQueue = deque()
            while queue:
                r, c = queue.popleft()
                if (r, c) in visited:
                    continue
                visited.add((r, c))
                for nr, nc in [(r - 1, c), (r + 1, c), (r, c - 1), (r, c + 1)]:
                    if 0 <= nr < row and 0 <= nc < col and (nr, nc) not in visited:
                        if matrix[nr][nc] == color:
                            queue.append((nr, nc))
                        else:
                            nextQueue.append((nr, nc))
            queue = nextQueue
            color ^= 1
        return res - 1


print(Solution().solve(matrix=[[0, 0, 0, 1], [1, 1, 1, 1], [0, 0, 0, 0]]))

# First set the color of (0, 0) to 1 and then to 0
