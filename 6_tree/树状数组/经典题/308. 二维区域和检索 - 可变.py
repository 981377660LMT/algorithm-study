from typing import List
from BIT import BIT4, BIT5

# 1 <= m, n <= 200
# 最多调用1e4 次 sumRegion 和 update 方法


class NumMatrix:
    def __init__(self, matrix: List[List[int]]):
        self.matrix = matrix
        self.row = len(matrix)
        self.col = len(matrix[0])
        self.bit = BIT4(self.row, self.col)
        for r in range(self.row):
            for c in range(self.col):
                self.bit.add(r, c, matrix[r][c])

    def update(self, row: int, col: int, val: int) -> None:
        """更新 matrix[row][col] 的值到 val 。"""
        delta = val - self.matrix[row][col]
        self.matrix[row][col] = val
        self.bit.add(row, col, delta)

    def sumRegion(self, row1: int, col1: int, row2: int, col2: int) -> int:
        """返回矩阵 matrix 中指定矩形区域元素的 和 ，该区域由 左上角 (row1, col1) 和 右下角 (row2, col2) 界定。"""
        return self.bit.queryRange(row1, col1, row2, col2)


class NumMatrix2:
    def __init__(self, matrix: List[List[int]]):
        self.matrix = matrix
        self.row = len(matrix)
        self.col = len(matrix[0])
        self.bit = BIT5(self.row, self.col)
        for r in range(self.row):
            for c in range(self.col):
                self.bit.addRange(r, c, r, c, matrix[r][c])

    def update(self, row: int, col: int, val: int) -> None:
        """更新 matrix[row][col] 的值到 val 。"""
        delta = val - self.matrix[row][col]
        self.matrix[row][col] = val
        self.bit.addRange(row, col, row, col, delta)

    def sumRegion(self, row1: int, col1: int, row2: int, col2: int) -> int:
        """返回矩阵 matrix 中指定矩形区域元素的 和 ，该区域由 左上角 (row1, col1) 和 右下角 (row2, col2) 界定。"""
        return self.bit.queryRange(row1, col1, row2, col2)
