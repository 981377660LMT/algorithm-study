# 交换列后的最大1子矩阵


class Solution:
    def solve(self, matrix):
        row, col = len(matrix), len(matrix[0])

        # 统计直方图的高度
        copy = [[0] * col for _ in range(row)]
        for c in range(col):
            copy[0][c] = matrix[0][c]
            for r in range(1, row):
                if matrix[r][c] == 1:
                    copy[r][c] = copy[r - 1][c] + 1
                else:
                    copy[r][c] = 0

        # 行排序，用桶排序
        for r in range(row):
            bucket = [0] * (row + 1)  # 桶排序
            for c in range(col):
                bucket[copy[r][c]] += 1
            index = 0
            for num in range(row, -1, -1):
                for _ in range(bucket[num]):
                    copy[r][index] = num
                    index += 1

        # 找到最大面积的1矩阵
        res = 0
        for r in range(row):
            for c in range(col):
                res = max(res, (c + 1) * copy[r][c])
        return res


print(Solution().solve(matrix=[[0, 0, 1], [1, 1, 1], [1, 0, 1]]))
print(Solution().solve(matrix=[[1], [0], [1]]))
# 4
# We can rearrange the columns to:

# [[0, 1, 0],
#  [1, 1, 1],
#  [1, 1, 0]]
# And then take the bottom 2 x 2 submatrix with all 1s

# matrix = [
#   [0, 1, 0, 1, 0],
#   [0, 1, 0, 1, 1],
#   [1, 1, 0, 1, 0]
# ]

# hist = [
#   [0, 1, 0, 1, 0],
#   [0, 2, 0, 2, 1],
#   [1, 3, 0, 3, 0]
# ]

# row-sorted hist = [
#   [1, 1, 0, 0, 0],
#   [2, 2, 1, 0, 0],
#   [3, 3, 1, 0, 0]
# ]

# area ending at every cell = [
#   [1, 2, 0, 0, 0],
#   [2, 4, 3, 0, 0],
#   [3, 6, 3, 0, 0],
# ]

# max area = 6
