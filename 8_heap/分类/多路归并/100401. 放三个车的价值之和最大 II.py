# 100401. 放三个车的价值之和最大 II
# https://leetcode.cn/problems/maximum-value-sum-by-placing-three-rooks-ii/solutions/2884123/top-3-omn-by-han3000-2hym/
# 给你一个 m x n 的二维整数数组 board ，它表示一个国际象棋棋盘，其中 board[i][j] 表示格子 (i, j) 的 价值 。
# 处于 同一行 或者 同一列 车会互相 攻击 。你需要在棋盘上放三个车，确保它们两两之间都 无法互相攻击 。
# 请你返回满足上述条件下，三个车所在格子 值 之和 最大 为多少。

# !每行每列只需要top3

from heapq import nlargest
from itertools import chain, combinations
from typing import List

INF = int(1e18)


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maximumValueSum(self, board: List[List[int]], k=3) -> int:
        ROW, COL = len(board), len(board[0])
        rows = [nlargest(k, [(board[i][j], i, j) for j in range(COL)]) for i in range(ROW)]
        cols = [nlargest(k, [(board[i][j], i, j) for i in range(ROW)]) for j in range(COL)]
        s = nlargest(k * k, set(chain(*rows)) & set(chain(*cols)))
        res = -INF
        for cand in combinations(s, k):
            if len(set(x[1] for x in cand)) == k and len(set(x[2] for x in cand)) == k:
                res = max2(res, sum(x[0] for x in cand))
        return res
