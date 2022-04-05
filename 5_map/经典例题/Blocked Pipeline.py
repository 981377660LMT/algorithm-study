# 2x(2*n)的矩阵 问 每次从左到右是否联通
# n<=1e5

# 用两个点表示一组blocked关系


class Solution:
    def solve(self, n, requests) -> int:
        res = 0
        grid = [[0] * 2 * n for _ in range(2)]
        blocked = set()  # 用两个点表示一组blocked关系

        for row, col, type in requests:
            grid[row][col] = type
            otherRow = row ^ 1

            if type == 1:
                for nextCol in (col - 1, col, col + 1):
                    if 0 <= nextCol < 2 * n and grid[otherRow][nextCol] == 1:
                        blocked.add((row, col, otherRow, nextCol))
            else:
                for nextCol in (col - 1, col, col + 1):
                    blocked.discard((row, col, otherRow, nextCol))
                    blocked.discard((otherRow, nextCol, row, col))

            if not blocked:
                res += 1

        return res


print(Solution().solve(n=4, requests=[[0, 2, 1], [1, 3, 1], [1, 3, 0]]))
# requests contains [row, col, type] meaning that m[row][col] becomes blocked
# if type = 1 and it becomes unblocked if type = 0.
