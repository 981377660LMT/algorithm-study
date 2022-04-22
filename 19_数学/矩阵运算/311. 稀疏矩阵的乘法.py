# https://www.cnblogs.com/grandyang/p/5282959.html
# 1.构造更为有效的数据结构
# 因为零值没有太多的意义，所以我们可以忽略零值，并且仅需要存储或操作稀疏矩阵中数据或非零值。有多种数据结构可用于有效地构造稀疏矩阵，下面列出了三个常见的例子。

# Dictionary of Keys：使用字典，其中行和列索引映射到值。
# List of Lists：矩阵的每一行都存储为一个列表，每个子列表包含列索引和值。
# Coordinate List ：每个元组都存储一个元组列表，其中包含行索引、列索引和值。

# 此处存为列表

from typing import List

# 矩阵乘法
class Solution:
    def multiply(self, mat1: List[List[int]], mat2: List[List[int]]) -> List[List[int]]:
        if len(mat1) == 0 or len(mat2) == 0:
            return []

        res = [[0] * len(mat2[0]) for _ in range(len(mat1))]
        lis1 = self.make_list(mat1)
        lis2 = self.make_list(mat2)

        for v1, r1, c1 in lis1:
            for v2, r2, c2 in lis2:
                if c1 == r2:
                    res[r1][c2] += v1 * v2
        return res

    @staticmethod
    def make_list(mat: List[List[int]]):
        lis = []
        for r in range(len(mat)):
            for c in range(len(mat[0])):
                if mat[r][c] != 0:
                    lis.append([mat[r][c], r, c])
        return lis


print(Solution().multiply([[1, 0, 0], [-1, 0, 3]], [[7, 0, 0], [0, 0, 0], [0, 0, 1]]))
