# k ≤ 20 where k is the number of unfilled squares.
# nxn 的矩阵，0表示还没有放置的棋子
# 判断每行每列是否能放置数字1-n


from collections import defaultdict
from itertools import product


class Solution:
    def solve(self, matrix) -> bool:
        n = len(matrix)
        rows, cols = defaultdict(set), defaultdict(set)
        for r, c in product(range(n), repeat=2):
            if matrix[r][c] in rows[r] or matrix[r][c] in cols[c]:
                return False
            if matrix[r][c] != 0:
                rows[r].add(matrix[r][c])
                cols[c].add(matrix[r][c])

        def bt(x: int, y: int) -> bool:
            if x == n:
                return True
            res = False
            if matrix[x][y] != 0:
                rows[x].add(matrix[x][y])
                cols[y].add(matrix[x][y])
                if y == n - 1:
                    res = res or bt(x + 1, 0)
                else:
                    res = res or bt(x, y + 1)
            else:
                for select in range(1, n + 1):
                    if select not in rows[x] and select not in cols[y]:
                        matrix[x][y] = select
                        rows[x].add(select)
                        cols[y].add(select)
                        if y == n - 1:
                            res = res or bt(x + 1, 0)
                        else:
                            res = res or bt(x, y + 1)
                        matrix[x][y] = 0
                        rows[x].remove(select)
                        cols[y].remove(select)
            return res

        return bt(0, 0)


# matrix = [
#     [1, 2, 3],
#     [2, 0, 1],
#     [0, 0, 2]
# ]
print(Solution().solve(matrix=[[1, 2, 3], [2, 0, 1], [0, 0, 2]]))
