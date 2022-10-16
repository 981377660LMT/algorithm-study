# 矩阵上下反转
from typing import List


def flipTopToBottom(matrix: List[List[int]]) -> List[List[int]]:
    """矩阵上下翻转"""
    newMatrix = [row[:] for row in matrix]
    ROW = len(matrix)
    for i in range(ROW // 2):
        newMatrix[i], newMatrix[~i] = newMatrix[~i], newMatrix[i]
    return newMatrix


if __name__ == "__main__":
    matrix = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
    print(flipTopToBottom(matrix))
