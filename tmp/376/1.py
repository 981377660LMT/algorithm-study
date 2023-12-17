from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的二维整数矩阵 grid，大小为 n * n ，其中的值在 [1, n2] 范围内。除了 a 出现 两次，b 缺失 之外，每个整数都 恰好出现一次 。

# 任务是找出重复的数字a 和缺失的数字 b 。


# 返回一个下标从 0 开始、长度为 2 的整数数组 ans ，其中 ans[0] 等于 a ，ans[1] 等于 b 。
class Solution:
    def findMissingAndRepeatedValues(self, grid: List[List[int]]) -> List[int]:
        counter = Counter()
        for i in range(len(grid)):
            for j in range(len(grid[0])):
                counter[grid[i][j]] += 1
        a, b = 0, 0
        for i in range(1, len(grid) * len(grid[0]) + 1):
            if counter[i] == 2:
                a = i
            if counter[i] == 0:
                b = i
        return [a, b]
