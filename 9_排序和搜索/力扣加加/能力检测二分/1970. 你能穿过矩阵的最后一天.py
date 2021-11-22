from collections import deque
from typing import List

# 你想知道从矩阵最 上面 一行走到最 下面 一行，且只经过陆地格子的 最后一天 是哪一天
class Solution:
    def latestDayToCross(self, row: int, col: int, cells: List[List[int]]) -> int:
        # 由于起点可能有多个（第一行的所有陆地），因此使用多源 BFS 复杂度会更好
        def can(mid: int) -> bool:
            visited = set()
            queue = deque([(0, i) for i in range(col)])
            for x, y in cells[:mid]:
                visited.add((x - 1, y - 1))
            while queue:
                x, y = queue.popleft()
                if (x, y) in visited:
                    continue
                visited.add((x, y))
                if x == row - 1:
                    return True
                for dx, dy in [(1, 0), (-1, 0), (0, 1), (0, -1)]:
                    if 0 <= x + dx < row and 0 <= y + dy < col:
                        queue.append((x + dx, y + dy))
            return False

        # 最右能力二分
        l = 0
        r = row * col
        while l <= r:
            mid = (l + r) >> 1
            if can(mid):
                l = mid + 1
            else:
                r = mid - 1
        return r


print(
    Solution().latestDayToCross(
        3, 3, [[1, 2], [2, 1], [3, 3], [2, 2], [1, 1], [1, 3], [2, 3], [3, 2], [3, 1]]
    )
)
