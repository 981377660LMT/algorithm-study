from typing import List
from collections import deque


# 黑子的邻居是一片连续白子，端点是黑子
# 白子翻转成黑子以后，黑子的位置可能作为新的端点，放入候选区
class Solution:
    def flipChess(self, chessboard: List[str]) -> int:
        dirs = [(-1, -1), (-1, 0), (-1, 1), (0, -1), (0, 1), (1, -1), (1, 0), (1, 1)]
        Row, Col = len(chessboard), len(chessboard[0])

        def bfs(row: int, col: int) -> int:
            mat = [list(row) for row in chessboard]
            res = 0
            queue = deque([(row, col)])
            while queue:
                cr, cc = queue.popleft()
                for di in range(8):
                    dr, dc = dirs[di]
                    nr = cr + dr
                    nc = cc + dc
                    if 0 <= nr < Row and 0 <= nc < Col:
                        if mat[nr][nc] == '.':  # 没有旗子
                            continue
                        elif mat[nr][nc] == 'X':  # 也是黑色，就无法延伸
                            continue
                        else:
                            while 0 <= nr < Row and 0 <= nc < Col and mat[nr][nc] == 'O':
                                nr += dr
                                nc += dc
                            # ----另一端是黑子
                            if 0 <= nr < Row and 0 <= nc < Col and mat[nr][nc] == 'X':
                                # ----把中间的白色，翻过来
                                nnr = cr + dr
                                nnc = cc + dc
                                while 0 <= nnr < Row and 0 <= nnc < Col and mat[nnr][nnc] == 'O':
                                    mat[nnr][nnc] = 'X'  # 翻转成黑色
                                    res += 1  # 翻转白子的个数，+1
                                    queue.append((nnr, nnc))
                                    nnr += dr
                                    nnc += dc
            return res

        res = 0
        for r in range(Row):
            for c in range(Col):
                if chessboard[r][c] == '.':
                    cur = bfs(r, c)
                    res = max(cur, res)
        return res


print(Solution().flipChess(["....X.", "....X.", "XOOO..", "......", "......"]))
