from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)

# 返回满足 Ri 行和 Cj 列相等的行列对 (Ri, Cj) 的数目。
# 如果行和列以相同的顺序包含相同的元素（即相等的数组），则认为二者是相等的。

# 总结:元组作为状态(哈希值)


class Solution:
    def equalPairs(self, grid: List[List[int]]) -> int:
        rowCounter = defaultdict(int)
        colCounter = defaultdict(int)
        for row in grid:
            state = tuple(row)
            rowCounter[state] += 1
        for col in zip(*grid):
            state = tuple(col)
            colCounter[state] += 1

        res = 0
        for c, v in rowCounter.items():
            res += v * colCounter[c]
        return res
