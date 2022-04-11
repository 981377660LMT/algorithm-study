# Flipped Matrix Prequel
# 反转一行后再反转一列 求最大1的个数

# 1 ≤ n, m ≤ 250
class Solution:
    def solve(self, matrix):
        row, col = len(matrix), len(matrix[0])
        rowCounter = [0] * row
        colCounter = [0] * col
        for i in range(row):
            for j in range(col):
                rowCounter[i] += matrix[i][j]
                colCounter[j] += matrix[i][j]

        sum_ = sum(rowCounter)
        res = 0
        for i in range(row):
            for j in range(col):
                # How many ones and zeros flip if flipping @ (r, c) ?
                ones = rowCounter[i] + colCounter[j] - 2 * matrix[i][j]  # 少了的1
                zeros = row + col - 2 - ones  # 多出来的0
                cand = sum_ - ones + zeros
                res = max(res, cand)
        return res


print(Solution().solve(matrix=[[1, 0, 1], [0, 1, 0], [1, 0, 0]]))
