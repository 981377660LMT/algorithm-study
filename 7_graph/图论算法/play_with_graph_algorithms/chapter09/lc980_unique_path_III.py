from typing import List

class Solution:

    DIRECTIONS = [(0, 1), (1, 0), (0, -1), (-1, 0)]

    def uniquePathsIII(self, grid: List[List[int]]) -> int:
        # 遍历所有能走的点一次
        # 有多少种走法
        self._grid = grid
        self._R = len(grid)
        self._C = len(grid[0])
        self._visited = [[False] * self._C for _ in range(self._R)]
        self._left = self._R * self._C

        for i in range(self._R):
            for j in range(self._C):
                if self._grid[i][j] == 1:
                    self._start = i * self._C + j
                    self._grid[i][j] = 0
                elif self._grid[i][j] == 2:
                    self._end = i * self._C + j
                    self._grid[i][j] = 0
                elif self._grid[i][j] == -1:
                    self._left -= 1

        return self._dfs(self._start, self._left)

    def _dfs(self, v, left):
        x = v // self._C
        y = v % self._C
        self._visited[x][y] = True
        left -= 1

        if left == 0 and v == self._end:
            self._visited[x][y] = False
            return 1

        res = 0
        for dx, dy in self.DIRECTIONS:
            nextx = x + dx
            nexty = y + dy
            if self._in_area(nextx, nexty) and self._grid[nextx][nexty] == 0 and not self._visited[nextx][nexty]:
                res += self._dfs(nextx * self._C + nexty, left)
        
        self._visited[x][y] = False
        return res

    def _in_area(self, x, y):
        return x >= 0 and x < self._R and y >= 0 and y < self._C


class SolutionV2:

    DIRECTIONS = [(0, 1), (1, 0), (0, -1), (-1, 0)]

    def uniquePathsIII(self, grid: List[List[int]]) -> int:
        # 遍历所有能走的点一次
        # 有多少种走法
        self._grid = grid
        self._R = len(grid)
        self._C = len(grid[0])
        self._left = self._R * self._C
        # 第一个维度是visited有多少种状态的组合
        # 第二个维度是visited这个点有多少种选择
        self._memo = [
            [-1] * (self._R * self._C)
            for _ in range(1 << self._R * self._C)
        ]

        for i in range(self._R):
            for j in range(self._C):
                if self._grid[i][j] == 1:
                    self._start = i * self._C + j
                    self._grid[i][j] = 0
                elif self._grid[i][j] == 2:
                    self._end = i * self._C + j
                    self._grid[i][j] = 0
                elif self._grid[i][j] == -1:
                    self._left -= 1

        visited = 0
        return self._dfs(visited, self._start, self._left)

    def _dfs(self, visited, v, left):
        if self._memo[visited][v] != -1:
            return self._memo[visited][v]
        visited += (1 << v)
        left -= 1

        if left == 0 and v == self._end:
            self._memo[visited][v] = 1
            return 1

        x = v // self._C
        y = v % self._C
        res = 0
        for dx, dy in self.DIRECTIONS:
            nextx = x + dx
            nexty = y + dy
            next_ = nextx * self._C + nexty
            if self._in_area(nextx, nexty) and self._grid[nextx][nexty] == 0 and (visited & (1 << next_) == 0):
                res += self._dfs(visited, next_, left)

        self._memo[visited][v] = res
        return res

    def _in_area(self, x, y):
        return x >= 0 and x < self._R and y >= 0 and y < self._C


if __name__ == '__main__':
    sol = Solution()
    grid = [[1,0,0,0],[0,0,0,0],[0,0,2,-1]]
    print(sol.uniquePathsIII(grid))

    sol = SolutionV2()
    grid = [[1,0,0,0],[0,0,0,0],[0,0,2,-1]]
    print(sol.uniquePathsIII(grid))