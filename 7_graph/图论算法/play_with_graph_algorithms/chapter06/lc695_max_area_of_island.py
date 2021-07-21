class Solution1:

    DIRECTIONS = [(0, 1), (1, 0), (0, -1), (-1, 0)]

    def max_area_of_island(self, grid):
        if not grid or not grid[0]:
            return 0
        self._R, self._C = len(grid), len(grid[0])
        self._grid = grid
        self._G = self._construct_graph()
        self._visited = [False] * len(self._G)

        res = 0
        for v in range(len(self._G)):
            x = v // self._C
            y = v % self._C
            if not self._visited[v] and self._grid[x][y] == 1:
                res = max(res, self._dfs(v))
        return res

    def _construct_graph(self):
        g = [set() for _ in range(self._R * self._C)]

        for v in range(len(g)):
            x = v // self._C
            y = v % self._C
            if self._grid[x][y] == 1:
                for dx, dy in self.DIRECTIONS:
                    nextx = x + dx
                    nexty = y + dy
                    if self._in_area(nextx, nexty) and self._grid[nextx][nexty] == 1:
                        next_ = nextx * self._C + nexty
                        g[v].add(next_)
                        g[next_].add(v)
        return g

    def _in_area(self, x, y):
        return 0 <= x < self._R and 0 <= y < self._C

    def _dfs(self, v):
        self._visited[v] = True
        res = 1
        for w in self._G[v]:
            if not self._visited[w]:
                res += self._dfs(w)
        return res


class Solution2:

    DIRECTIONS = [(0, 1), (1, 0), (0, -1), (-1, 0)]

    def max_area_of_island(self, grid):
        if not grid or not grid[0]:
            return 0
        self._R, self._C = len(grid), len(grid[0])
        self._grid = grid
        self._visited = [[False] * self._C for _ in range(self._R)]
        res = 0
        for i in range(self._R):
            for j in range(self._C):
                if not self._visited[i][j] and self._grid[i][j] == 1:
                    res = max(res, self._dfs(i, j))
        return res

    def _in_area(self, x, y):
        return 0 <= x < self._R and 0 <= y < self._C

    def _dfs(self, x, y):
        self._visited[x][y] = True
        res = 1
        for dx, dy in self.DIRECTIONS:
            nextx, nexty = x + dx, y + dy
            if self._in_area(nextx, nexty) and not self._visited[nextx][nexty] and self._grid[nextx][nexty]:
                res += self._dfs(nextx, nexty)
        return res


class Solution3:

    DIRECTIONS = [(0, 1), (1, 0), (0, -1), (-1, 0)]

    def max_area_of_island(self, grid):
        max_area = 0
        if not grid or not grid[0]:
            return max_area
        m, n = len(grid), len(grid[0])

        for i in range(m):
            for j in range(n):
                if grid[i][j] == 1:
                    max_area = max(max_area, self._dfs(grid, i, j, 0))
        return max_area

    # return the max area starting from i, j
    def _dfs(self, grid, i, j, curr):
        grid[i][j] = 2
        curr += 1

        for di, dj in self.DIRECTIONS:
            newi, newj = i + di, j + dj
            if not 0 <= newi < len(grid) or not 0 <= newj < len(grid[0]):
                continue
            if grid[newi][newj] != 1:
                continue
            curr = max(curr, self._dfs(grid, newi, newj, curr))

        return curr


class UF:

    def __init__(self, n):
        self._parent = [i for i in range(n)]
        self._sz = [1] * n

    def _find(self, p):
        if p != self._parent[p]:
            self._parent[p] = self._find(self._parent[p])
        return self._parent[p]

    def is_connected(self, p, q):
        return self._find(p) == self._find(q)

    def union_elements(self, p, q):
        p_root = self._find(p)
        q_root = self._find(q)
        if p_root == q_root:
            return
        self._parent[p_root] = q_root
        self._sz[q_root] += self._sz[p_root]

    def size(self, p):
        return self._sz[self._find(p)]


class Solution4:

    DIRECTIONS = [(0, 1), (1, 0), (0, -1), (-1, 0)]

    def max_area_of_island(self, grid):
        if not grid or not grid[0]:
            return 0
        
        R, C = len(grid), len(grid[0])
        uf = UF(R * C)

        for v in range(R * C):
            x = v // C
            y = v % C
            if grid[x][y] != 1:
                continue
            for dx, dy in self.DIRECTIONS:
                nextx, nexty = x + dx, y + dy
                if not 0 <= nextx < R or not 0 <= nexty < C:
                    continue
                if grid[nextx][nexty] != 1:
                    continue
                next_ = nextx * C + nexty
                uf.union_elements(v, next_)

        res = 0
        for v in range(R * C):
            x = v // C
            y = v % C
            if grid[x][y] == 1:
                res = max(res, uf.size(v))
        return res



if __name__ == '__main__':
    data = [
        [1, 1, 0, 0, 0],
        [1, 1, 0, 0, 0],
        [0, 0, 0, 1, 1],
        [0, 0, 0, 1, 1],
    ]

    sol1 = Solution1()
    print(sol1.max_area_of_island(data))

    sol2 = Solution2()
    print(sol2.max_area_of_island(data))

    sol3 = Solution3()
    print(sol3.max_area_of_island(data))

    data = [
        [1, 1, 0, 0, 0],
        [1, 1, 0, 0, 0],
        [0, 0, 0, 1, 1],
        [0, 0, 0, 1, 1],
    ]
    sol4 = Solution4()
    print(sol4.max_area_of_island(data))