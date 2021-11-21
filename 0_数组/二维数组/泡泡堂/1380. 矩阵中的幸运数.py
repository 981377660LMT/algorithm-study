# 按 任意 顺序返回矩阵中的所有幸运数。

# 幸运数是指矩阵中满足同时下列两个条件的元素：

# 在同一行的所有元素中最小
# 在同一列的所有元素中最大
# 矩阵中的数字 各不相同
from typing import List


class Solution:
    def luckyNumbers(self, matrix: List[List[int]]) -> List[int]:
        rowmin = [min(i) for i in matrix]
        colmin = [max(i) for i in zip(*matrix)]
        return list(set(rowmin) & set(colmin))
