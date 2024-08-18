from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 m x n 的二维整数数组 board ，它表示一个国际象棋棋盘，其中 board[i][j] 表示格子 (i, j) 的 价值 。

# 处于 同一行 或者 同一列 车会互相 攻击 。你需要在棋盘上放三个车，确保它们两两之间都 无法互相攻击 。


# 请你返回满足上述条件下，三个车所在格子 值 之和 最大 为多少。
class Solution:
    def maximumValueSum(self, board: List[List[int]]) -> int:
        ROW, COL = len(board), len(board[0])
        order = [(i, j) for i in range(ROW) for j in range(COL)]
        order.sort(key=lambda x: board[x[0]][x[1]], reverse=True)

        return
