from typing import Deque, List, Tuple
from collections import deque

# "@" 是起点，（"a", "b", ...）代表钥匙，（"A", "B", ...）代表锁。
# "." 代表一个空房间， "#" 代表一堵墙，
# 返回获取所有钥匙所需要的移动的最少次数。如果无法获取所有钥匙，返回 -1 。
# 钥匙的数目范围是 [1, 6]，每个钥匙都对应一个不同的字母
# m,n <=30

# 还是 visited + queue

# 思路:
# 1. 预处理，找起点和钥匙个数，每把🔑赋值一个id
# 2. bfs 队列存 row,col,visitedKeys
# 3. 使用三维数组来表示visited (每个点的每个状态是否访问过)


class Solution:
    def shortestPathAllKeys(self, grid: List[str]) -> int:
        if not any(grid):
            return -1
        row, col = len(grid), len(grid[0])
        sr, sc = 0, 0
        key_id = dict()
        key_count = 0

        for r in range(row):
            for c in range(col):
                if grid[r][c] == '@':  # 找起点
                    sr, sc = r, c
                elif grid[r][c].islower():  # 找钥匙
                    key_id[grid[r][c]] = key_count
                    key_count += 1

        target = (1 << key_count) - 1  # 终点end == 所有钥匙都到手了

        queue: Deque[Tuple[int, int, int, int]] = deque([(sr, sc, 0, 0)])
        # 多了一维状态
        visited = [[[False for _ in range(1 << key_count)] for _ in range(col)] for _ in range(row)]

        while queue:
            # print(Q)
            r, c, state, dist = queue.popleft()

            if state == target:
                return dist

            for nr, nc in [(r - 1, c), (r + 1, c), (r, c - 1), (r, c + 1)]:
                if 0 <= nr < row and 0 <= nc < col and grid[nr][nc] != '#':  # 不是墙
                    nv = grid[nr][nc]
                    if nv.isupper():  # 是lock
                        k = nv.lower()
                        if (state >> key_id[k]) & 1:  # 如果有key
                            if not visited[nr][nc][state]:
                                visited[nr][nc][state] = True
                                queue.append((nr, nc, state, dist + 1))  # 继续往下走

                    elif nv.islower():  # 是key
                        next_state = state | (1 << key_id[nv])
                        if not visited[nr][nc][next_state]:
                            visited[nr][nc][next_state] = True
                            queue.append((nr, nc, next_state, dist + 1))

                    elif nv in ('.', '@'):
                        if not visited[nr][nc][state]:
                            visited[nr][nc][state] = True
                            queue.append((nr, nc, state, dist + 1))

        return -1


print(Solution().shortestPathAllKeys(["@.a.#", "###.#", "b.A.B"]))
# 输出：8
