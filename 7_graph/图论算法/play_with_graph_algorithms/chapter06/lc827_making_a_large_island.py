class UF:

    def __init__(self, n):
        self._parent = [i for i in range(n)]
        self._sz = [1] * n
        self._max = 1

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
        self._max = max(self._max, self._sz[q_root])

    def size(self, p):
        return self._sz[self._find(p)]

    def max(self):
        return self._max


class Solution:
    def largest_island(self, grid):
        if not grid or not grid[0]:
            return 0

        m, n = len(grid), len(grid[0])
        uf = UF(m * n)
        directions = [(0, 1), (1, 0), (0, -1), (-1, 0)]
        for v in range(m * n):
            i, j = v // m, v % n
            if grid[i][j] != 1:
                continue
            for di, dj in directions:
                newi, newj = i + di, j + dj
                if not 0 <= newi < m or not 0 <= newj < n:
                    continue
                if grid[newi][newj] != 1:
                    continue
                next_ = newi * n + newj
                uf.union_elements(v, next_)

        res = -2 ** 31
        for i in range(m):
            for j in range(n):
                if grid[i][j] != 0:
                    continue
                surrounded_parents = set()
                for di, dj in directions:
                    newi, newj = i + di, j + dj
                    if not 0 <= newi < m or not 0 <= newj < n:
                        continue
                    if grid[newi][newj] == 0:
                        continue
                    surrounded_parents.add(uf._find(newi * n + newj))
                temp = 0
                for each in surrounded_parents:
                    temp += uf.size(each)
                res = max(res, temp + 1)

        return res if res != -2 ** 31 else uf.max()
                    

if __name__ == '__main__':
    sol = Solution()
    data = [[1, 1], [1, 0]]
    print(sol.largestIsland(data))