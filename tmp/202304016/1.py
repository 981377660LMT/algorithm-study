from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个大小为 m x n 的二进制矩阵 mat ，请你找出包含最多 1 的行的下标（从 0 开始）以及这一行中 1 的数目。

# 如果有多行包含最多的 1 ，只需要选择 行下标最小 的那一行。


# 返回一个由行下标和该行中 1 的数量组成的数组。
class Solution:
    def rowAndMaximumOnes(self, mat: List[List[int]]) -> List[int]:
        ones = [row.count(1) for row in mat]
        max_ = max(ones)
        maxRow = ones.index(max_)
        return [maxRow, max_]
