from typing import List

# 转置口诀:zip星负一
class Solution:
    def rotate(self, matrix: List[List[int]]) -> None:
        """
        Do not return anything, modify matrix in-place instead.
        """
        matrix[:] = zip(*matrix[::-1])
        # zip(*matrix)表示转置矩阵


a = [
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9],
]
print(*a[::-1])
