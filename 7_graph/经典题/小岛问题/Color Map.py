# Return the minimum number of operations required so that every cell has the same color.

# 统计连通分量的颜色种类，需要染的就是sum(counter.values())-max(counter.values())
from collections import Counter


class Solution:
    def solve(self, A):
        if not A:
            return 0
        R, C = len(A), len(A[0])

        counter = Counter()
        visited = set()

        def dfs(r, c, color):
            for nr, nc in ((r - 1, c), (r + 1, c), (r, c - 1), (r, c + 1)):
                if 0 <= nr < R and 0 <= nc < C and A[nr][nc] == color:
                    if (nr, nc) not in visited:
                        visited.add((nr, nc))
                        dfs(nr, nc, color)

        for r in range(R):
            for c in range(C):
                if (r, c) not in visited:
                    visited.add((r, c))
                    dfs(r, c, A[r][c])
                    counter[A[r][c]] += 1

        # count = {1: 3, 2: 1, 3: 1}
        vals = counter.values()  # vals = [3, 1, 1]
        return sum(vals) - max(vals)
