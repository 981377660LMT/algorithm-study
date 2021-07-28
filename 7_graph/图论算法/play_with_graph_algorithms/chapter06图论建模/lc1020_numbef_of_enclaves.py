class Solution:

    DIRECTIONS = [(0, 1), (1, 0), (0, -1), (-1, 0)]

    def num_of_enclaves(self, A):
        if not A or not A[0]:
            return 0

        self._m, self._n = len(A), len(A[0])
        visited = set()
        for i in range(self._m):
            if A[i][0] == 1:
                self._dfs_initial_boundary(A, i, 0, visited)
            if A[i][self._n - 1] == 1:
                self._dfs_initial_boundary(A, i, self._n - 1, visited)
        for j in range(self._n):
            if A[0][j] == 1:
                self._dfs_initial_boundary(A, 0, j, visited)
            if A[self._m - 1][j] == 1:
                self._dfs_initial_boundary(A, self._m - 1, j, visited)

        res = 0
        for i in range(self._m):
            for j in range(self._n):
                if A[i][j] == 1:
                    self._dfs(A, i, j, visited)
                    res += 1

        for i in range(self._m):
            for j in range(self._n):
                if A[i][j] == 2:
                    A[i][j] = 1

        return res

    def _dfs_initial_boundary(self, A, i, j, visited):
        A[i][j] = 2
        visited.add((i, j))
        for di, dj in self.DIRECTIONS:
            newi, newj = i + di, j + dj
            if not 0 <= newi < self._m or not 0 <= newj < self._n:
                continue
            if A[newi][newj] == 0 or A[newi][newj] == 2:
                continue
            if (newi, newj) in visited:
                continue
            self._dfs_initial_boundary(A, newi, newj, visited)

    def _dfs(self, A, i, j, visited):
        visited.add((i, j))
        for di, dj in self.DIRECTIONS:
            newi, newj = i + di, j + dj
            if not 0 <= newi < self._m or not 0 <= newj < self._n:
                continue
            if A[newi][newj] != 1:
                continue
            if (newi, newj) in visited:
                continue
            self._dfs(A, newi, newj, visited)


if __name__ == '__main__':
    sol = Solution()
    # data = [[0, 0, 0, 0], [1, 0, 1, 0], [0, 1, 1, 0], [0, 0, 0, 0]]
    data = [[0,0,1,1,1,0,1,1,1,0,1],[1,1,1,1,0,1,0,1,1,0,0],[0,1,0,1,1,0,0,0,0,1,0],[1,0,1,1,1,1,1,0,0,0,1],[0,0,1,0,1,1,0,0,1,0,0],[1,0,0,1,1,1,0,0,0,1,1],[0,1,0,1,1,0,0,0,1,0,0],[0,1,1,0,1,0,1,1,1,0,0],[1,1,0,1,1,1,0,0,0,0,0],[1,0,1,1,0,0,0,1,0,0,1]]
    print(sol.num_of_enclaves(data))
