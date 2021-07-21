class Solution:

    # 主要思路就是先沿着四边向中心扩展，将O变成#
    # 然后把所有剩余的O（此时这些O肯定不与边界相连）直接变成X
    # 最后将#变回O即可
    def solve(self, board):
        if not board or not board[0]:
            return
        
        m, n = len(board), len(board[0])
        for i in range(m):
            for j in range(n):
                if (i in (0, m - 1) or j in (0, n - 1)) and board[i][j] == 'O':
                    self._dfs(board, i, j)
        
        for i in range(m):
            for j in range(n):
                if board[i][j] == 'O':
                    board[i][j] = 'X'
                if board[i][j] == '#':
                    board[i][j] = 'O'
    
    def _dfs(self, board, i, j):
        board[i][j] = '#'

        m, n = len(board), len(board[0])
        if i > 0 and board[i - 1][j] == 'O':
            self._dfs(board, i - 1, j)
        if j > 0 and board[i][j - 1] == 'O':
            self._dfs(board, i, j - 1)
        if i < m - 1 and board[i + 1][j] == 'O':
            self._dfs(board, i + 1, j)
        if j < n - 1 and board[i][j + 1] == 'O':
            self._dfs(board, i, j + 1)