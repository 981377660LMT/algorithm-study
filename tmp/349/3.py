from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 n 和一个下标从 0 开始的 二维数组 queries ，其中 queries[i] = [typei, indexi, vali] 。

# 一开始，给你一个下标从 0 开始的 n x n 矩阵，所有元素均为 0 。每一个查询，你需要执行以下操作之一：


# 如果 typei == 0 ，将第 indexi 行的元素全部修改为 vali ，覆盖任何之前的值。
# 如果 typei == 1 ，将第 indexi 列的元素全部修改为 vali ，覆盖任何之前的值。
# 请你执行完所有查询以后，返回矩阵中所有整数的和。


# 倒过来
class Solution:
    def matrixSumQueries(self, n: int, queries: List[List[int]]) -> int:
        ...
