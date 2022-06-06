from typing import List
from BIT import BIT4


class NumMatrix:
    def __init__(self, matrix: List[List[int]]):
        self.matrix = matrix
        self.row = len(matrix)
        self.col = len(matrix[0])
        self.bit = BIT4(self.row, self.col)
        for r in range(self.row):
            for c in range(self.col):
                self.bit.update(r + 1, c + 1, matrix[r][c])

    def update(self, row: int, col: int, val: int) -> None:
        """更新 matrix[row][col] 的值到 val 。"""
        delta = val - self.matrix[row][col]
        self.matrix[row][col] = val
        self.bit.update(row + 1, col + 1, delta)

    def sumRegion(self, row1: int, col1: int, row2: int, col2: int) -> int:
        """返回矩阵 matrix 中指定矩形区域元素的 和 ，该区域由 左上角 (row1, col1) 和 右下角 (row2, col2) 界定。"""
        return self.bit.sumRange(row1 + 1, col1 + 1, row2 + 1, col2 + 1)

