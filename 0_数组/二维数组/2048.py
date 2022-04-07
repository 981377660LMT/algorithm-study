# 2048游戏
# 4x4的矩阵
from typing import List, Literal

Matrix = List[List[int]]
Direction = Literal['left', 'right', 'up', 'down']


def rotate(matrix: Matrix, times=1) -> Matrix:
    """顺时针旋转矩阵90度`times`次"""
    assert times >= 1
    res = [list(col[::-1]) for col in zip(*matrix)]
    for _ in range(times - 1):
        res = [list(col[::-1]) for col in zip(*res)]
    return res


class Solution:
    def solve(self, board: Matrix, direction: Direction) -> Matrix:
        # 统一成left
        if direction == 'right':
            board = rotate(board, 2)
        elif direction == 'up':
            board = rotate(board, 3)
        elif direction == 'down':
            board = rotate(board, 1)

        for r in range(4):
            row = [num for num in board[r] if num != 0]
            for c in range(3):
                if c + 1 < len(row) and row[c] == row[c + 1]:
                    row[c] *= 2
                    row.pop(c + 1)
            while len(row) < 4:
                row.append(0)
            board[r] = row

        if direction == 'right':
            board = rotate(board, 2)
        elif direction == 'up':
            board = rotate(board, 1)
        elif direction == 'down':
            board = rotate(board, 3)

        return board


print(
    Solution().solve(
        board=[[2, 0, 0, 2], [2, 2, 2, 2], [0, 4, 2, 2], [2, 2, 2, 0]], direction='left'
    )
)
# [
#     [4, 0, 0, 0],
#     [4, 4, 0, 0],
#     [4, 4, 0, 0],
#     [4, 2, 0, 0]
# ]
