from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个由 正整数 组成、大小为 m x n 的矩阵 grid。你可以从矩阵中的任一单元格移动到另一个位于正下方或正右侧的任意单元格（不必相邻）。从值为 c1 的单元格移动到值为 c2 的单元格的得分为 c2 - c1 。

# 你可以从 任一 单元格开始，并且必须至少移动一次。


# 返回你能得到的 最大 总得分。
class Solution:
    def maxScore(self, grid: List[List[int]]) -> int:
        ...
