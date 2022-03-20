from numpy.linalg import matrix_power
import numpy as np


class Solution:
    def climbStairs(self, n: int) -> int:
        a = np.array([[1, 1], [1, 0]])
        b = np.array([[2], [1]])
        return int(np.matmul(matrix_power(a, n - 1), b)[1])

