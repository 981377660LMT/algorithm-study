from typing import List

# 技巧:全部处理成字符串lines(lines = board + list(cols) + list(diag))


class Solution:
    def validTicTacToe(self, board: List[str]) -> bool:

        """ 
        1. `X` 的数量总是等于`O`的数量或者比它多1 （保证先手顺序）
        2. 只能有一个胜者，或者没有胜者而棋盘是满的
        3. `X`获胜时，`X`的棋子数量比`O`多1，且`O`未获胜
        4. `O`获胜时，`O`的棋子数量和`X`相同，且`X`未获胜
        """

        xNum = sum(sum(x == 'X' for x in line) for line in board)
        oNum = sum(sum(x == 'O' for x in line) for line in board)

        if (xNum != oNum) and (xNum - oNum != 1):
            return False

        # 这一段来自checkIO上x-o-referee最简洁的解法
        # zip一个二维数组=>提取列
        cols = map(''.join, zip(*board))
        diag = map(''.join, zip(*[(r[i], r[2 - i]) for i, r in enumerate(board)]))
        lines = board + list(cols) + list(diag)

        oWin = True if 'OOO' in lines else False
        xWin = True if 'XXX' in lines else False

        if oWin and xWin:
            return False

        if oWin and (xNum != oNum):
            return False

        if xWin and (xNum - oNum != 1):
            return False

        return True
