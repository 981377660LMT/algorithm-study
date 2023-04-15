from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的 m x n 整数矩阵 grid 。矩阵中某一列的宽度是这一列数字的最大 字符串长度 。

# 比方说，如果 grid = [[-10], [3], [12]] ，那么唯一一列的宽度是 3 ，因为 -10 的字符串长度为 3 。
# 请你返回一个大小为 n 的整数数组 ans ，其中 ans[i] 是第 i 列的宽度。


# 一个有 len 个数位的整数 x ，如果是非负数，那么 字符串长度 为 len ，否则为 len + 1 。
class Solution:
    def findColumnWidth(self, grid: List[List[int]]) -> List[int]:
        res = []
        for i in range(len(grid[0])):
            res.append(max(len(str(grid[j][i])) for j in range(len(grid))))
        return res
