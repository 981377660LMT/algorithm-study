from collections import defaultdict
from typing import List


class Solution:
    def solveSudoku(self, board: List[List[str]]) -> None:
        """
        Do not return anything, modify board in-place instead.

        原地解数独
        """

        def bt(row: int, col: int) -> bool:
            if col == n:
                return bt(row + 1, 0)
            if row == n:
                return True
            if board[row][col] != ".":
                return bt(row, col + 1)

            for cur in "123456789":
                if (
                    cur not in rowVisited[row]
                    and cur not in colVisited[col]
                    and cur not in blockVisited[(row // 3, col // 3)]
                ):
                    board[row][col] = cur
                    rowVisited[row].add(cur)
                    colVisited[col].add(cur)
                    blockVisited[(row // 3, col // 3)].add(cur)
                    if bt(row, col + 1):
                        return True
                    board[row][col] = "."
                    rowVisited[row].remove(cur)
                    colVisited[col].remove(cur)
                    blockVisited[(row // 3, col // 3)].remove(cur)
            return False

        n = len(board)
        rowVisited, colVisited, blockVisited = defaultdict(set), defaultdict(set), defaultdict(set)
        for r in range(n):
            for c in range(n):
                if board[r][c] != ".":
                    rowVisited[r].add(board[r][c])
                    colVisited[c].add(board[r][c])
                    blockVisited[(r // 3, c // 3)].add(board[r][c])

        bt(0, 0)
