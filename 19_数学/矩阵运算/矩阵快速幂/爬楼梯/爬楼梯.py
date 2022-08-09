import numpy as np


def matqpow(base: np.ndarray, exp: int, mod: int) -> np.ndarray:
    """np矩阵快速幂"""

    base = base.copy()
    res = np.eye(*base.shape, dtype=np.uint64)

    while exp:
        if exp & 1:
            res = res @ base
        exp //= 2
        base = base @ base
    return res


class Solution:
    def climbStairs(self, n: int) -> int:
        """矩阵乘法"""
        if n <= 3:
            return n
        res = np.array([[2], [1]])
        trans = np.array([[1, 1], [1, 0]])
        tmp = matqpow(trans, n - 2, int(1e9 + 7))
        res = tmp @ res
        return int(res[0][0])


print(Solution().climbStairs(3))
