class Solution:

    _DIRS = [
        (-1, -1), (-1, 0), (-1, 1), (0, -1),
        (0, 1), (1, -1), (1, 0), (1, 1),
    ]

    def update_board(self, board, click):
        if not board or not board[0]:
            return board
        
        m, n = len(board), len(board[0])
        row, col = click
        
        if board[row][col] == 'M':
            board[row][col] = 'X'
        else:
            mine_counts = 0
            for di, dj in self._DIRS:
                newi, newj = row + di, col + dj
                if not 0 <= newi < m or not 0 <= newj < n:
                    continue
                if board[newi][newj] == 'M':
                    mine_counts += 1
            if mine_counts > 0:
                board[row][col] = str(mine_counts)
            else:
                board[row][col] = 'B'
                for di, dj in self._DIRS:
                    newi, newj = row + di, col + dj
                    if not 0 <= newi < m or not 0 <= newj < n:
                        continue
                    if board[newi][newj] == 'E':
                        self.updateBoard(board, [newi, newj])
                
        return board
