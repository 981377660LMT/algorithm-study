from typing import List

# dfs
class Solution1:
    def numIslands(self, grid: List[List[str]]) -> int:
        if not grid:
            return 0
        row = len(grid)
        col = len(grid[0])
        cnt = 0

        def dfs(i, j):
            grid[i][j] = "0"
            for x, y in [[-1, 0], [1, 0], [0, -1], [0, 1]]:
                tmp_i = i + x
                tmp_j = j + y
                if 0 <= tmp_i < row and 0 <= tmp_j < col and grid[tmp_i][tmp_j] == "1":
                    dfs(tmp_i, tmp_j)

        for i in range(row):
            for j in range(col):
                if grid[i][j] == "1":
                    dfs(i, j)
                    cnt += 1
        return cnt


# bfs
class Solution2:
    def numIslands(self, grid: List[List[str]]) -> int:
        from collections import deque

        if not grid:
            return 0
        row = len(grid)
        col = len(grid[0])
        cnt = 0

        def bfs(i, j):
            queue = deque()
            queue.appendleft((i, j))
            grid[i][j] = "0"
            while queue:
                i, j = queue.pop()
                for x, y in [[-1, 0], [1, 0], [0, -1], [0, 1]]:
                    tmp_i = i + x
                    tmp_j = j + y
                    if 0 <= tmp_i < row and 0 <= tmp_j < col and grid[tmp_i][tmp_j] == "1":
                        grid[tmp_i][tmp_j] = "0"
                        queue.appendleft((tmp_i, tmp_j))

        for i in range(row):
            for j in range(col):
                if grid[i][j] == "1":
                    bfs(i, j)
                    cnt += 1
        return cnt


# 并查集：看有几个根
class Solution3:
    def numIslands(self, grid: List[List[str]]) -> int:
        f = {}

        def find(x):
            f.setdefault(x, x)
            if f[x] != x:
                f[x] = find(f[x])
            return f[x]

        def union(x, y):
            f[find(x)] = find(y)

        if not grid:
            return 0
        row = len(grid)
        col = len(grid[0])

        for i in range(row):
            for j in range(col):
                if grid[i][j] == "1":
                    for x, y in [[-1, 0], [0, -1]]:
                        tmp_i = i + x
                        tmp_j = j + y
                        if 0 <= tmp_i < row and 0 <= tmp_j < col and grid[tmp_i][tmp_j] == "1":
                            union(tmp_i * row + tmp_j, i * row + j)
        # print(f)
        res = set()
        for i in range(row):
            for j in range(col):
                if grid[i][j] == "1":
                    res.add(find((i * row + j)))
        return len(res)

