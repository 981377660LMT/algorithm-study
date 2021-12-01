from typing import List
from collections import defaultdict

# 矩阵 mat 有 6 行 3 列，从 mat[2][0] 开始的 `矩阵对角线` 将会经过 mat[2][0]、mat[3][1] 和 mat[4][2] 。
# ，请你将同一条 矩阵对角线 上的元素按升序排序后，返回排好序的矩阵。
class Solution:
    def diagonalSort(self, mat: List[List[int]]) -> List[List[int]]:
        if not any(mat):
            return []

        m, n = len(mat), len(mat[0])
        dic = defaultdict(list)
        for i in range(m):
            for j in range(n):
                dic[i - j].append(mat[i][j])

        for key in dic:
            dic[key].sort(reverse=True)

        for i in range(m):
            for j in range(n):
                mat[i][j] = dic[i - j].pop()

        return mat


print(Solution().diagonalSort(mat=[[3, 3, 1, 1], [2, 2, 1, 2], [1, 1, 1, 2]]))
