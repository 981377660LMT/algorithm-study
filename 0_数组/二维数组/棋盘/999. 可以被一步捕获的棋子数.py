from typing import List

# #  8 x 8 的棋盘上有一个白色的车（Rook）
# # 车 Rook  R
# # 象 Bishop B
# # 黑卒 pawn p
# # 空地 .
# # 你现在可以控制车移动一次，请你统计有多少敌方的卒处于你的捕获范围内（即，可以被一步捕获的棋子数）。


class Solution:
    def numRookCaptures(self, board: List[List[str]]) -> int:
        # 找到车的位置
        rook_row, rook_col = next((i, j) for i in range(8) for j in range(8) if board[i][j] == 'R')
        row_chess = ''.join(board[rook_row][c] for c in range(8) if board[rook_row][c] != '.')
        col_chess = ''.join(board[r][rook_col] for r in range(8) if board[r][rook_col] != '.')
        # print(row_chess, col_chess)
        # print(list(s for s in (row_chess, col_chess)))
        # 'Rp':右/下 'pR':左/上
        return sum('Rp' in line for line in (row_chess, col_chess)) + sum(
            'pR' in line for line in (row_chess, col_chess)
        )


print(
    Solution().numRookCaptures(
        [
            [".", ".", ".", ".", ".", ".", ".", "."],
            [".", ".", ".", "p", ".", ".", ".", "."],
            [".", ".", ".", "p", ".", ".", ".", "."],
            ["p", "p", ".", "R", ".", "p", "B", "."],
            [".", ".", ".", ".", ".", ".", ".", "."],
            [".", ".", ".", "B", ".", ".", ".", "."],
            [".", ".", ".", "p", ".", ".", ".", "."],
            [".", ".", ".", ".", ".", ".", ".", "."],
        ]
    )
)

