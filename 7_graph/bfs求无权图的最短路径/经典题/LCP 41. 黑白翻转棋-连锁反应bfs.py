from typing import List
from collections import deque


# 黑子的邻居是一片连续白子，端点是黑子
# 白子翻转成黑子以后，黑子的位置可能作为新的端点，放入候选区

# LCP 41. 黑白翻转棋-连锁反应用bfs 类似水纹扩散 炸弹引爆

DIRS = [(-1, -1), (-1, 0), (-1, 1), (0, -1), (0, 1), (1, -1), (1, 0), (1, 1)]


class Solution:
    def flipChess(self, chessboard: List[str]) -> int:
        def bfs(rowStart: int, colStart: int) -> int:
            """在(row,col)处放一个黑子，返回反转白子的数量"""
            grid = [list(row) for row in chessboard]
            res = 0
            queue = deque([(rowStart, colStart)])

            while queue:
                r, c = queue.popleft()
                for dr, dc in DIRS:
                    nr, nc = r + dr, c + dc
                    if 0 <= nr < ROW and 0 <= nc < COL:
                        if grid[nr][nc] == '.':  # 没有旗子
                            continue
                        elif grid[nr][nc] == 'X':  # 也是黑色，就无法延伸
                            continue
                        else:
                            while 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 'O':
                                nr += dr
                                nc += dc

                            # ----另一端是黑子
                            if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 'X':
                                # ----把中间的白色，翻过来
                                nnr = r + dr
                                nnc = c + dc
                                while 0 <= nnr < ROW and 0 <= nnc < COL and grid[nnr][nnc] == 'O':
                                    grid[nnr][nnc] = 'X'  # 翻转成黑色
                                    res += 1  # 翻转白子的个数，+1
                                    queue.append((nnr, nnc))
                                    nnr += dr
                                    nnc += dc
            return res

        ROW, COL = len(chessboard), len(chessboard[0])

        res = 0
        for r in range(ROW):
            for c in range(COL):
                if chessboard[r][c] == '.':
                    cur = bfs(r, c)
                    res = max(cur, res)

        return res


print(Solution().flipChess(["....X.", "....X.", "XOOO..", "......", "......"]))
