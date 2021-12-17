from typing import List

# 这个问题是实现一个简单的消除算法。
# 官方答案的递归解法很优雅 特殊用途的网络 Ad-Hoc
# https://leetcode-cn.com/problems/candy-crush/solution/fen-sui-tang-guo-by-leetcode/

# 1. 如果有`三个及以上水平或者垂直`相连的同种糖果，`同一时间`将它们粉碎，即将这些位置变成空的。
# 2. 在同时粉碎掉这些糖果之后，如果有一个空的位置上方还有糖果，那么上方的糖果就会下落直到碰到下方的糖果或者底部

# 思路：
# 1.垂直，水平扫描整个表格粉碎糖果；对要粉碎的糖果做标记而不是当场粉碎(否则会影响其他结果)；
# 使用滑动窗口或者itertools.groupby来判断连续
# 2.掉落糖果模拟(栈或者逆序遍历列，读写指针)；
# 3.是否继续要使用一个标志位判断
class Solution:
    def candyCrush(self, board: List[List[int]]) -> List[List[int]]:
        row, col = len(board), len(board[0])
        todo = False

        # 粉碎标记
        for r in range(row):
            for c in range(col - 2):
                if abs(board[r][c]) == abs(board[r][c + 1]) == abs(board[r][c + 2]) != 0:
                    board[r][c] = board[r][c + 1] = board[r][c + 2] = -abs(board[r][c])
                    todo = True

        for r in range(row - 2):
            for c in range(col):
                if abs(board[r][c]) == abs(board[r + 1][c]) == abs(board[r + 2][c]) != 0:
                    board[r][c] = board[r + 1][c] = board[r + 2][c] = -abs(board[r][c])
                    todo = True

        # 下落，逆序读写指针
        for c in range(col):
            write = row - 1
            for read in range(row - 1, -1, -1):
                if board[read][c] > 0:
                    board[write][c] = board[read][c]
                    write -= 1
            for r in range(write, -1, -1):
                board[r][c] = 0

        return board if not todo else self.candyCrush(board)


print(
    Solution().candyCrush(
        board=[
            [110, 5, 112, 113, 114],
            [210, 211, 5, 213, 214],
            [310, 311, 3, 313, 314],
            [410, 411, 412, 5, 414],
            [5, 1, 512, 3, 3],
            [610, 4, 1, 613, 614],
            [710, 1, 2, 713, 714],
            [810, 1, 2, 1, 1],
            [1, 1, 2, 2, 2],
            [4, 1, 4, 4, 1014],
        ]
    )
)

