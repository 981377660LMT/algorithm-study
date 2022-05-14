from typing import List
import collections

# 2 <= n <= 100

# 注意状态是(x1,y1,x2,y2)
# 移动时需要分🐍水平还是竖直讨论


class Solution:
    def minimumMoves(self, grid: List[List[int]]) -> int:
        n = len(grid)
        start = (0, 0, 0, 1)
        end = (n - 1, n - 2, n - 1, n - 1)
        if grid[0][0] == 1 or grid[0][1] == 1 or grid[n - 1][n - 2] == 1 or grid[n - 1][n - 1] == 1:
            return -1

        queue = collections.deque()
        visited = set()
        queue.append(start)
        visited.add(start)
        step = 0
        while queue:
            cur_len = len(queue)
            for _ in range(cur_len):
                r1, c1, r2, c2 = queue.popleft()
                if (r1, c1, r2, c2) == end:
                    return step

                # ------------------------------------ 蛇身水平
                if r1 == r2:
                    # ------------- 右侧是空的，往右走一步
                    if c2 + 1 < n and grid[r2][c2 + 1] != 1 and (r2, c2, r2, c2 + 1) not in visited:
                        queue.append((r2, c2, r2, c2 + 1))
                        visited.add((r2, c2, r2, c2 + 1))
                    # ------------- 下侧都是空的
                    if r1 + 1 < n and grid[r1 + 1][c1] != 1 and grid[r1 + 1][c2] != 1:
                        # ---- 下移
                        if (r1 + 1, c1, r2 + 1, c2) not in visited:
                            queue.append((r1 + 1, c1, r2 + 1, c2))
                            visited.add((r1 + 1, c1, r2 + 1, c2))
                        # ---- 顺时针
                        if (r1, c1, r1 + 1, c1) not in visited:
                            queue.append((r1, c1, r1 + 1, c1))
                            visited.add((r1, c1, r1 + 1, c1))

                # ------------------------------------ 蛇身竖直
                if c1 == c2:
                    # ------------- 下侧是空的，往下一步
                    if r2 + 1 < n and grid[r2 + 1][c2] != 1 and (r2, c2, r2 + 1, c2) not in visited:
                        queue.append((r2, c2, r2 + 1, c2))
                        visited.add((r2, c2, r2 + 1, c2))
                    # ------------ 右侧2点都是空的
                    if c1 + 1 < n and grid[r1][c1 + 1] != 1 and grid[r2][c2 + 1] != 1:
                        # ---- 右移
                        if (r1, c1 + 1, r2, c2 + 1) not in visited:
                            queue.append((r1, c1 + 1, r2, c2 + 1))
                            visited.add((r1, c1 + 1, r2, c2 + 1))
                        # ---- 逆时针
                        if (r1, c1, r1, c1 + 1) not in visited:
                            queue.append((r1, c1, r1, c1 + 1))
                            visited.add((r1, c1, r1, c1 + 1))
            step += 1

        return -1

