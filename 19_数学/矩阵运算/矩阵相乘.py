from typing import List
import numpy as np


class Solution:
    def multiply1(self, mat1: List[List[int]], mat2: List[List[int]]) -> List[List[int]]:
        return [[sum(a * b for a, b in zip(row, col)) for col in zip(*mat2)] for row in mat1]

    def multiply2(self, mat1: List[List[int]], mat2: List[List[int]]) -> List[List[int]]:
        return np.dot(mat1, mat2).tolist()

