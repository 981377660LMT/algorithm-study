# k ≤ 20 where k is the number of unfilled squares.
# nxn 的矩阵，0表示还没有放置的棋子
# 判断每行每列是否能放置数字1-n
# 类似于解数独

from collections import defaultdict
from itertools import product


class Solution:
    def solve(self, matrix) -> bool:
        def bt(row: int, col: int) -> bool:
            if col == n:
                return bt(row + 1, 0)
            if row == n:
                return True
            if matrix[row][col] != 0:
                return bt(row, col + 1)

            for select in range(1, n + 1):
                if select not in rowVisited[row] and select not in colVisited[col]:
                    matrix[row][col] = select
                    rowVisited[row].add(select)
                    colVisited[col].add(select)
                    if bt(row, col + 1):
                        return True
                    matrix[row][col] = 0
                    rowVisited[row].remove(select)
                    colVisited[col].remove(select)

            return False

        n = len(matrix)
        rowVisited, colVisited = defaultdict(set), defaultdict(set)
        for r, c in product(range(n), repeat=2):
            if matrix[r][c] in rowVisited[r] or matrix[r][c] in colVisited[c]:
                return False
            if matrix[r][c] != 0:
                rowVisited[r].add(matrix[r][c])
                colVisited[c].add(matrix[r][c])

        return bt(0, 0)


# matrix = [
#     [1, 2, 3],
#     [2, 0, 1],
#     [0, 0, 2]
# ]
print(Solution().solve(matrix=[[1, 2, 3], [2, 0, 1], [0, 0, 2]]))
