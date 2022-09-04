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

DIR4 = ((0, 1), (0, -1), (1, 0), (-1, 0))


class Solution:
    def shortestPathAllKeys(self, grid: List[str]) -> int:
        ROW, COL = len(grid), len(grid[0])
        sr, sc = 0, 0
        keyId = dict()
        keyCount = 0

        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == "@":
                    sr, sc = r, c
                elif grid[r][c].islower():
                    keyId[grid[r][c]] = keyCount
                    keyCount += 1

        target = (1 << keyCount) - 1  # 所有钥匙都被拿到了
        queue: Deque[Tuple[int, int, int, int]] = deque([(sr, sc, 0, 0)])
        visited = [[[False] * (1 << keyCount) for _ in range(COL)] for _ in range(ROW)]  # 多了一维状态

        while queue:
            curRow, curCol, keyState, curDist = queue.popleft()

            if keyState == target:
                return curDist

            for dr, dc in DIR4:
                nr, nc = curRow + dr, curCol + dc
                if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] != "#":  # 不是墙
                    nv = grid[nr][nc]
                    if nv.isupper():  # 是lock
                        k = nv.lower()
                        if (keyState >> keyId[k]) & 1 and not visited[nr][nc][keyState]:  # 如果有key
                            visited[nr][nc][keyState] = True
                            queue.append((nr, nc, keyState, curDist + 1))  # 继续往下走

                    elif nv.islower():  # 是key
                        nextState = keyState | (1 << keyId[nv])
                        if not visited[nr][nc][nextState]:
                            visited[nr][nc][nextState] = True
                            queue.append((nr, nc, nextState, curDist + 1))

                    elif nv in (".", "@"):
                        if not visited[nr][nc][keyState]:
                            visited[nr][nc][keyState] = True
                            queue.append((nr, nc, keyState, curDist + 1))

        return -1


print(Solution().shortestPathAllKeys(["@.a.#", "###.#", "b.A.B"]))
# 输出：8
