from typing import List

# 每次移动，你能任意交换两列或是两行的位置。
# board 是方阵，且行列数的范围是[2, 30]。
# 输出将这个矩阵变为 “棋盘” 所需的最小移动次数。
class Solution:
    def movesToChessboard(self, board: List[List[int]]) -> int:
        ...


print(Solution().movesToChessboard(board=[[0, 1, 1, 0], [0, 1, 1, 0], [1, 0, 0, 1], [1, 0, 0, 1]]))


# 没太懂、
