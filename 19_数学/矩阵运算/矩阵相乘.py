from typing import List
import numpy as np
from scipy.sparse import csr_matrix


class Solution:
    def multiply1(self, mat1: List[List[int]], mat2: List[List[int]]) -> List[List[int]]:
        return [[sum(a * b for a, b in zip(row, col)) for col in zip(*mat2)] for row in mat1]

    def multiply2(self, mat1: List[List[int]], mat2: List[List[int]]) -> List[List[int]]:
        return np.dot(mat1, mat2).tolist()

    def multiply3(self, mat1: List[List[int]], mat2: List[List[int]]) -> List[List[int]]:
        """稀疏矩阵的乘法
        
        https://www.runoob.com/scipy/scipy-sparse-matrix.html
        """
        m1, m2 = csr_matrix(mat1), csr_matrix(mat2)
        res = m1 * m2
        return res.todense().tolist()  # 转为稠密矩阵


print(Solution().multiply3(mat1=[[1, 0, 0], [-1, 0, 3]], mat2=[[7, 0, 0], [0, 0, 0], [0, 0, 1]]))
