from typing import List


class Solution:
    def updateBoard(self, board: List[List[str]], click: List[int]) -> List[List[str]]:
        i, j, m, n = *click, len(board), len(board[0])

        def dfs(i, j):
            if board[i][j] == 'E':
                neis = [
                    (x, y)
                    for x, y in (
                        (i - 1, j - 1),
                        (i - 1, j),
                        (i - 1, j + 1),
                        (i, j - 1),
                        (i, j + 1),
                        (i + 1, j - 1),
                        (i + 1, j),
                        (i + 1, j + 1),
                    )
                    if 0 <= x < m and 0 <= y < n
                ]
                cnt = sum(board[x][y] == 'M' for x, y in neis)
                # 如果一个没有相邻地雷的空方块（'E'）被挖出，修改它为（'B'），并且所有和其相邻的未挖出方块都应该被递归地揭露。
                if not cnt:
                    board[i][j] = 'B'
                    for x, y in neis:
                        dfs(x, y)
                # 如果一个至少与一个地雷相邻的空方块（'E'）被挖出，修改它为数字（'1'到'8'），表示相邻地雷的数量。
                else:
                    board[i][j] = str(cnt)

        if board[i][j] == 'M':
            board[i][j] = 'X'
        else:
            dfs(i, j)
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
