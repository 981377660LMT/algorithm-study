from typing import List

DIR8 = [(0, 1), (1, 0), (0, -1), (-1, 0), (1, 1), (1, -1), (-1, 1), (-1, -1)]


class Solution:
    def updateBoard(self, board: List[List[str]], click: List[int]) -> List[List[str]]:
        def getNexts(r: int, c: int):
            for dr, dc in DIR8:
                nr, nc = r + dr, c + dc
                if 0 <= nr < ROW and 0 <= nc < COL:
                    yield nr, nc

        def dfs(r: int, c: int) -> None:
            if board[r][c] == 'E':
                nexts = list(getNexts(r, c))
                count = sum(board[x][y] == 'M' for x, y in nexts)
                # 如果一个没有相邻地雷的空方块（'E'）被挖出，修改它为（'B'），并且所有和其相邻的未挖出方块都应该被递归地揭露。
                if not count:
                    board[r][c] = 'B'
                    for nr, nc in nexts:
                        dfs(nr, nc)
                # 如果一个至少与一个地雷相邻的空方块（'E'）被挖出，修改它为数字（'1'到'8'），表示相邻地雷的数量。
                else:
                    board[r][c] = str(count)

        sr, sc, ROW, COL = *click, len(board), len(board[0])

        if board[sr][sc] == 'M':
            board[sr][sc] = 'X'
        else:
            dfs(sr, sc)
        return board


# console.log(
#   updateBoard(
#     [
#       ['E', 'E', 'E', 'E', 'E'],
#       ['E', 'E', 'M', 'E', 'E'],
#       ['E', 'E', 'E', 'E', 'E'],
#       ['E', 'E', 'E', 'E', 'E'],
#     ],
#     [3, 0]
#   )
# )
# // 输出:

# // [['B', '1', 'E', '1', 'B'],
# //  ['B', '1', 'M', '1', 'B'],
# //  ['B', '1', '1', '1', 'B'],
# //  ['B', 'B', 'B', 'B', 'B']]

# 如果一个地雷（'M'）被挖出，游戏就结束了- 把它改为 'X'。
# 如果一个没有相邻地雷的空方块（'E'）被挖出，修改它为（'B'），并且所有和其相邻的未挖出方块都应该被递归地揭露。
# 如果一个至少与一个地雷相邻的空方块（'E'）被挖出，修改它为数字（'1'到'8'），表示相邻地雷的数量。
# 如果在此次点击中，若无更多方块可被揭露，则返回面板。
